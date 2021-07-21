set BASE_DIR=%~dp0
set Filepath=FilepathToSubstitution


if exist %Filepath% (
	set result=0
) else (
	set result=1
)

if %result% EQU 0 (
   rd /s /q %Filepath%
   if exist %Filepath% (
		set result2=0
   ) else (
		set result2=1
   )
) else (
	set result2=1
)
echo %result2%


