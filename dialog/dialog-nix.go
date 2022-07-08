//go:build (linux || unix)

package dialog

func find() {
    // "In older versions of Go, LookPath could return a path relative to the
    // current directory. As of Go 1.19, LookPath will instead return that path
    // along with an error satisfying errors.Is(err, ErrDot). See the package
    // documentation for more details."

    // In case imported by an older version of Go, we use "which"

    // Due to build constraints, we don't have to care about Windows
    // being insecure if it looks in the current directory for exec.Command

    // return (exec.Command("sh", "-c", "which "+cmd+" > /dev/null 2>&1").Run() == nil)

    // TODO TODO TODO
}
