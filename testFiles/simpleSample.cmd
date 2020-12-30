@echo off
setlocal EnableDelayedExpansion
goto process_start
;Important! Do not modify above this line!

;Begin profile definitions. While all fields are required, of particular importance is the format field, which includes a format description and extension separated by a |. This field is used as a filter when selecting files to convert in the GUI, therefore an incorrectly-formatted format field will result in an inoperable profile!

[PROFILE]
name="Sonic the Hedgehog"
format="CRI Movie 2 (*.usm)|*.usm"
example="hedge.usm"
description="Gotta go fast!"

;Begin menu definitions. MENU0 and MENU1 are reserved for the profile and target directory settings, respectively, therefore profile menus must begin at MENU2. As many menus as defined in the current profile are displayed in the GUI, however only a maximum of 10 arguments may be passed into batch code*. Therefore, menus are limited from 0-9.
;(*more arguments are possible using the "shift" command, but this is not directly supported by Vidsquish)

;Each menu definition is written as a default string value and one or more 'cases' numbered from 0-infinity. Each case is written as a string label paired with a value which will be passed into batch code. Label and value pairs are separated by a |. All values are passed into batch code as strings. Any unspecified values are treated as -1. The GUI will not allow proceeding until no values equal -1.

;Optionally, each menu definition may also include a text description which will cause a help icon to be displayed beside the menu. Clicking this help icon will display description text defined for the menu. If a description is not declared, a help icon will not be displayed.

[MENU2] ;Aspect ratio
case0="16:9 (HD)|16_9"
case1="21:9 (Ultra-wide)|21_9"
case2="32:9 (Super ultra-wide)|32_9"
default="Select aspect ratio"

[MENU3] ;Resolution
case0="Optimal (Higher quality)|high"
case1="Reduced (Higher performance)|low"
default="Select resolution profile"
description="You can pick your favorite resolution profile."

;Begin batch processing. All code following the BATCH section will be executed directly once for each file in the target directory, therefore no further .ini definitions should be declared beyond this point.

;Arguments are input in the same order as menus above, where %0 is reserved as the batch file itself and %1 is reserved as the current input file, including the full file path. As such, arguments begin with %2 corresponding to the value set by MENU2.

;Note that batch code is executed from the temporary "%APPDATA%\Vidsquish\lib" directory created at runtime.

[BATCH]
:process_start

REM Some comments here

echo "Hello world!"
