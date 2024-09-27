param($jabba_home, $jabba_symlink)

while ("True")
{
    $jabba_home = Read-Host 'input jabba home directory'
    if ($jabba_home)
    {
        if (Test-Path $jabba_home -PathType Leaf)
        {
            Write-Host "The input directory is is an existing file. reinput jabba home directory"
            continue
        }

        if (!(Test-Path $jabba_home -PathType Container))
        {
            New-Item -ItemType Directory -Force -Path $jabba_home
        }
        $jabba_home = (Resolve-Path $jabba_home).Path
        break
    }
    else
    {
        Write-Host "jabba home directory can not be null or empty."
    }
}


while ("True")
{
    $jabba_symlink = Read-Host 'input jabba symlink jdk directory'
    if ($jabba_symlink)
    {
        if (Test-Path $jabba_symlink -PathType Leaf)
        {
            Write-Host "The input directory is is an existing file. reinput jabba symlink jdk directory"
            continue
        }
        if ((Test-Path $jabba_symlink -PathType Container))
        {
            Write-Host "The current directory cannot be an existing one"
            continue
        }
        New-Item -ItemType Directory -Force -Path $jabba_symlink
        $jabba_symlink = (Resolve-Path $jabba_symlink).Path
        break
    }
    else
    {
        Write-Host "jabba symlink jdk directory can not be null or empty."
    }
}

Write-Host "set JABBA_HOME = $jabba_home"
Write-Host "set JABBA_SYMLINK = $jabba_symlink"

cmd.exe /C setx JABBA_HOME $jabba_home
cmd.exe /C setx JABBA_SYMLINK $jabba_symlink

Write-Host "Please Add current to Env Path: [%JABBA_HOME%] or [$jabba_home]"
Write-Host "Please Add current to Env Path: [%JABBA_HOME%] or [$jabba_home]"
Write-Host "Please Add current to Env Path: [%JABBA_HOME%] or [$jabba_home]"
