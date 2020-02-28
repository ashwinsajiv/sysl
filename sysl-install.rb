require 'fileutils'

Src = File.join('/', 'Users', ENV['USER'], 'Library', 'Caches', 'Homebrew', 'sysl-0.9.0', 'cmd', 'sysl')
Dir.chdir Src
system("go", "build")
Dest = File.join(ENV['HOME'],'go', 'bin')
FileUtils.cp(File.join(Src, 'sysl') ,Dest)
