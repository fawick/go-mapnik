cls

@echo off
setlocal EnableDelayedExpansion

If NOT DEFINED MAPNIK_SDK_PATH set MAPNIK_SDK_PATH=c:\mapnik-sdk

echo.
echo.Using Mapnik SDK from here: %MAPNIK_SDK_PATH%
echo.

If NOT EXIST %MAPNIK_SDK_PATH%\bin\mapnik-config.bat (
    echo.ERROR: Cannot find mapnik-config bat in %MAPNIK_SDK_PATH%\bin. Aborting...
    goto :eof
)

echo.
echo.Downloading C API from github
echo.

if not exist mapnik_c_api.cpp curl -LO https://raw.github.com/fawick/mapnik-c-api/master/mapnik_c_api.cpp
if not exist mapnik_c_api.h curl -LO https://raw.github.com/fawick/mapnik-c-api/master/mapnik_c_api.h

If DEFINED ProgramFiles(x86) Set BUILDTOOLS32BIT=%ProgramFiles(x86)%
If NOT DEFINED ProgramFiles(x86) Set BUILDTOOLS32BIT=%ProgramFiles%

call "%BUILDTOOLS32BIT%\Microsoft Visual Studio 10.0\VC\vcvarsall.bat" x86

echo.
echo.Compiling C API to shared library
echo.

set MAPNIK_C_DEFINES=-DUNICODE -DWIN32

set PATH=%PATH%;%MAPNIK_SDK_PATH%\bin

for /f "delims=" %%i in ('mapnik-config --prefix') do set MAPNIK_PATH=%%i
for /f "delims=" %%i in ('mapnik-config --cxxflags') do set MAPNIK_CXXFLAGS=%%i
for /f "delims=" %%i in ('mapnik-config --ldflags') do set MAPNIK_LDFLAGS=%%i
for /f "delims=" %%i in ('mapnik-config --dep-libs') do set MAPNIK_DEPLIBS=%%i
for /f "delims=" %%i in ('mapnik-config --libs') do set MAPNIK_LIBS=%%i
for /f "delims=" %%i in ('mapnik-config --defines') do set MAPNIK_DEFINES_LIST=%%i
for /f "delims=" %%i in ('mapnik-config --fonts') do set MAPNIK_FONTS_DIRECTORY=%%i
for /f "delims=" %%i in ('mapnik-config --input-plugins') do set MAPNIK_INPUT_PLUGINS_DIRECTORY=%%i


call :constructDefines %MAPNIK_DEFINES_LIST%

echo.LDFLAGS=%MAPNIK_LDFLAGS%
echo.LIBS=%MAPNIK_LIBS%
echo.DEFINES=%MAPNIK_C_DEFINES%
echo.FONTS=%MAPNIK_FONTS_DIRECTORY%
echo.PLUGINS=%MAPNIK_INPUT_PLUGINS_DIRECTORY%


> gen_import_windows.go (
    echo.package mapnik
	echo.
    echo.// #cgo LDFLAGS: %MAPNIK_SDK_PATH%/lib/mapnik_c_api.dll
	echo.import"C"
	echo.
	echo.const (
	echo.    fontPath = "%MAPNIK_FONTS_DIRECTORY%"
	echo.    pluginPath = "%MAPNIK_INPUT_PLUGINS_DIRECTORY%"
	echo.^)
)

echo.
echo.Compiling C API to shared library
echo.

if not exist mapnik_c_api.obj cl -c -nologo -Zm200 -Zc:wchar_t- -O2 %MAPNIK_CXXFLAGS% -W3 -w34100 -w34189 %MAPNIK_C_DEFINES% -I"%MAPNIK_SDK_PATH%\include" -I"%MAPNIK_SDK_PATH%\include\mapnik\agg" -I"." mapnik_c_api.cpp
::link /LIBPATH:%MAPNIK_LDFLAGS% %MAPNIK_LIBS% %MAPNIK_DEPLIBS% mapnik_c_api.obj /NOLOGO /DYNAMICBASE /NXCOMPAT /INCREMENTAL:NO /DLL /OUT:mapnik_c_api.dll 
link /LIBPATH:%MAPNIK_LDFLAGS% %MAPNIK_LIBS%  mapnik_c_api.obj /DLL /OUT:mapnik_c_api.dll 

IF NOT EXIST mapnik_c_api.dll (
    echo.
    echo.ERROR: Could not compile the DLL. Check for any error messages above.
    goto :eof
)

echo.
echo.Installing C API DLL to %MAPNIK_SDK_PATH%\lib
echo.

del mapnik_c_api.obj  mapnik_c_api.cpp  mapnik_c_api.lib  mapnik_c_api.exp
move /y mapnik_c_api.dll %MAPNIK_SDK_PATH%\lib

echo.
echo.Installing Go package
echo.

go install 

goto :eof

:constructDefines
	if "%1"=="" exit /b
	set MAPNIK_C_DEFINES=%MAPNIK_C_DEFINES% -D%1
	shift
	goto :constructDefines


