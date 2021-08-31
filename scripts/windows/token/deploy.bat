set BASE_DIR=%~dp0
set Filepath=%2
set Sourcepath=%BASE_DIR%%4
echo %BASE_DIR%
echo %Filepath%
echo %Sourcepath%

if exist %Filepath% (
	set result=1
) else (
	set result=0
)

if %result% EQU 1 (
    copy /y %Sourcepath% %Filepath%
	if exist %Filepath% (
		set result2=1
	) else (
		set result2=0
	)
) else (
	set result2=1
)
echo %result2%
