# bash completion for srn                               -*- shell-script -*-

__srn_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE:-} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__srn_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__srn_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__srn_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__srn_handle_go_custom_completion()
{
    __srn_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly srn allows handling aliases
    args=("${words[@]:1}")
    # Disable ActiveHelp which is not supported for bash completion v1
    requestComp="srn_ACTIVE_HELP=0 ${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __srn_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __srn_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __srn_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __srn_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __srn_debug "${FUNCNAME[0]}: the completions are: ${out}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __srn_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __srn_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __srn_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __srn_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subdir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out}")
        if [ -n "$subdir" ]; then
            __srn_debug "Listing directories in $subdir"
            __srn_handle_subdirs_in_dir_flag "$subdir"
        else
            __srn_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out}" -- "$cur")
    fi
}

__srn_handle_reply()
{
    __srn_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __srn_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION:-}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi

            if [[ -z "${flag_parsing_disabled}" ]]; then
                # If flag parsing is enabled, we have completed the flags and can return.
                # If flag parsing is disabled, we may not know all (or any) of the flags, so we fallthrough
                # to possibly call handle_go_custom_completion.
                return 0;
            fi
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __srn_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __srn_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        if declare -F __srn_custom_func >/dev/null; then
            # try command name qualified custom func
            __srn_custom_func
        else
            # otherwise fall back to unqualified for compatibility
            declare -F __custom_func >/dev/null && __custom_func
        fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__srn_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__srn_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__srn_handle_flag()
{
    __srn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue=""
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __srn_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __srn_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __srn_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __srn_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        __srn_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__srn_handle_noun()
{
    __srn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __srn_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __srn_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__srn_handle_command()
{
    __srn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_srn_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __srn_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__srn_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __srn_handle_reply
        return
    fi
    __srn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __srn_handle_flag
    elif __srn_contains_word "${words[c]}" "${commands[@]}"; then
        __srn_handle_command
    elif [[ $c -eq 0 ]]; then
        __srn_handle_command
    elif __srn_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __srn_handle_command
        else
            __srn_handle_noun
        fi
    else
        __srn_handle_noun
    fi
    __srn_handle_word
}

_srn_archive_unzip_password()
{
    last_command="srn_archive_unzip_password"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_archive_unzip()
{
    last_command="srn_archive_unzip"

    command_aliases=()

    commands=()
    commands+=("password")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_archive_zip_password()
{
    last_command="srn_archive_zip_password"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_archive_zip()
{
    last_command="srn_archive_zip"

    command_aliases=()

    commands=()
    commands+=("password")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_archive()
{
    last_command="srn_archive"

    command_aliases=()

    commands=()
    commands+=("unzip")
    commands+=("zip")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_completion()
{
    last_command="srn_completion"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("bash")
    must_have_one_noun+=("fish")
    must_have_one_noun+=("powershell")
    must_have_one_noun+=("zsh")
    noun_aliases=()
}

_srn_container_build()
{
    last_command="srn_container_build"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_delete()
{
    last_command="srn_container_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_images_delete()
{
    last_command="srn_container_images_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_images_list()
{
    last_command="srn_container_images_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_images()
{
    last_command="srn_container_images"

    command_aliases=()

    commands=()
    commands+=("delete")
    commands+=("list")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_ios()
{
    last_command="srn_container_ios"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--pair")
    flags+=("-p")
    local_nonpersistent_flags+=("--pair")
    local_nonpersistent_flags+=("-p")
    flags+=("--sidestore")
    flags+=("-s")
    local_nonpersistent_flags+=("--sidestore")
    local_nonpersistent_flags+=("-s")
    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_list()
{
    last_command="srn_container_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_shell()
{
    last_command="srn_container_shell"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ip")
    local_nonpersistent_flags+=("--ip")
    flags+=("--mount")
    flags+=("-m")
    local_nonpersistent_flags+=("--mount")
    local_nonpersistent_flags+=("-m")
    flags+=("--temp")
    flags+=("-t")
    local_nonpersistent_flags+=("--temp")
    local_nonpersistent_flags+=("-t")
    flags+=("--usb")
    flags+=("-u")
    local_nonpersistent_flags+=("--usb")
    local_nonpersistent_flags+=("-u")
    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container_silent()
{
    last_command="srn_container_silent"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ip")
    local_nonpersistent_flags+=("--ip")
    flags+=("--mount")
    flags+=("-m")
    local_nonpersistent_flags+=("--mount")
    local_nonpersistent_flags+=("-m")
    flags+=("--usb")
    flags+=("-u")
    local_nonpersistent_flags+=("--usb")
    local_nonpersistent_flags+=("-u")
    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_container()
{
    last_command="srn_container"

    command_aliases=()

    commands=()
    commands+=("build")
    commands+=("delete")
    commands+=("images")
    commands+=("ios")
    commands+=("list")
    commands+=("shell")
    commands+=("silent")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_dir_delete()
{
    last_command="srn_find_dir_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_dir()
{
    last_command="srn_find_dir"

    command_aliases=()

    commands=()
    commands+=("delete")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_file_delete()
{
    last_command="srn_find_file_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_file()
{
    last_command="srn_find_file"

    command_aliases=()

    commands=()
    commands+=("delete")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_word_delete()
{
    last_command="srn_find_word_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find_word()
{
    last_command="srn_find_word"

    command_aliases=()

    commands=()
    commands+=("delete")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_find()
{
    last_command="srn_find"

    command_aliases=()

    commands=()
    commands+=("dir")
    commands+=("file")
    commands+=("word")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch_create()
{
    last_command="srn_git_branch_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch_list()
{
    last_command="srn_git_branch_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch_local()
{
    last_command="srn_git_branch_local"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch_remote()
{
    last_command="srn_git_branch_remote"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch_switch()
{
    last_command="srn_git_branch_switch"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_branch()
{
    last_command="srn_git_branch"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("list")
    commands+=("local")
    commands+=("remote")
    commands+=("switch")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_changes()
{
    last_command="srn_git_changes"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_changes()
{
    last_command="srn_git_commit_changes"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_compare()
{
    last_command="srn_git_commit_compare"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_delete()
{
    last_command="srn_git_commit_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_list()
{
    last_command="srn_git_commit_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_push()
{
    last_command="srn_git_commit_push"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit_undo()
{
    last_command="srn_git_commit_undo"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_commit()
{
    last_command="srn_git_commit"

    command_aliases=()

    commands=()
    commands+=("changes")
    commands+=("compare")
    commands+=("delete")
    commands+=("list")
    commands+=("push")
    commands+=("undo")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_remote()
{
    last_command="srn_git_remote"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_stage()
{
    last_command="srn_git_stage"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_status()
{
    last_command="srn_git_status"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_sync()
{
    last_command="srn_git_sync"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_tag_create()
{
    last_command="srn_git_tag_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_tag_local()
{
    last_command="srn_git_tag_local"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_tag_remote()
{
    last_command="srn_git_tag_remote"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_tag_wipe()
{
    last_command="srn_git_tag_wipe"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_tag()
{
    last_command="srn_git_tag"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("local")
    commands+=("remote")
    commands+=("wipe")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_undo()
{
    last_command="srn_git_undo"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git_unstage()
{
    last_command="srn_git_unstage"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_git()
{
    last_command="srn_git"

    command_aliases=()

    commands=()
    commands+=("branch")
    commands+=("changes")
    commands+=("commit")
    commands+=("remote")
    commands+=("stage")
    commands+=("status")
    commands+=("sync")
    commands+=("tag")
    commands+=("undo")
    commands+=("unstage")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_help()
{
    last_command="srn_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_srn_music_convert_mp3()
{
    last_command="srn_music_convert_mp3"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_music_convert_playlist()
{
    last_command="srn_music_convert_playlist"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_music_convert()
{
    last_command="srn_music_convert"

    command_aliases=()

    commands=()
    commands+=("mp3")
    commands+=("playlist")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_music_download()
{
    last_command="srn_music_download"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_music()
{
    last_command="srn_music"

    command_aliases=()

    commands=()
    commands+=("convert")
    commands+=("download")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_clean()
{
    last_command="srn_nix_clean"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_home_build()
{
    last_command="srn_nix_home_build"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_home_delete()
{
    last_command="srn_nix_home_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_home_gen()
{
    last_command="srn_nix_home_gen"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_home()
{
    last_command="srn_nix_home"

    command_aliases=()

    commands=()
    commands+=("build")
    commands+=("delete")
    commands+=("gen")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_lint()
{
    last_command="srn_nix_lint"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_search()
{
    last_command="srn_nix_search"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_sys_build()
{
    last_command="srn_nix_sys_build"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_sys_delete()
{
    last_command="srn_nix_sys_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_sys_gen()
{
    last_command="srn_nix_sys_gen"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_sys()
{
    last_command="srn_nix_sys"

    command_aliases=()

    commands=()
    commands+=("build")
    commands+=("delete")
    commands+=("gen")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix_update()
{
    last_command="srn_nix_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_nix()
{
    last_command="srn_nix"

    command_aliases=()

    commands=()
    commands+=("clean")
    commands+=("home")
    commands+=("lint")
    commands+=("search")
    commands+=("sys")
    commands+=("update")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_todo()
{
    last_command="srn_todo"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_srn_root_command()
{
    last_command="srn"

    command_aliases=()

    commands=()
    commands+=("archive")
    commands+=("completion")
    commands+=("container")
    commands+=("find")
    commands+=("git")
    commands+=("help")
    commands+=("music")
    commands+=("nix")
    commands+=("todo")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_srn()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __srn_init_completion -n "=" || return
    fi

    local c=0
    local flag_parsing_disabled=
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("srn")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function=""
    local last_command=""
    local nouns=()
    local noun_aliases=()

    __srn_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_srn srn
else
    complete -o default -o nospace -F __start_srn srn
fi

# ex: ts=4 sw=4 et filetype=sh
