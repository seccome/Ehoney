set Filepath=%2
echo %Filepath%

if exist %Filepath% (
	set result=0
) else (
	set result=1
)

if %result% EQU 0 (
   del %Filepath%
   if exist %Filepath% (
		set result2=0
   ) else (
		set result2=1
   )
) else (
	set result2=1
)
echo %result2%