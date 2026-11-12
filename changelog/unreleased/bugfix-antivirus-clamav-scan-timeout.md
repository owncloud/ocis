Bugfix: Prevent virus scanning from getting stuck

If the antivirus scanner stopped responding, an upload's virus scan would
wait forever instead of failing. Enough stuck uploads could stop virus
scanning for everyone, blocking new uploads across the whole instance.

The virus scan now times out and fails instead of hanging indefinitely.
