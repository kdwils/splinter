Output vhs/merge.gif

Set Shell "bash"
Set FontSize 20
Set Width 1200
Set Height 600
Set WindowBar Colorful
Set CursorBlink false
Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10

Type "cat examples/split/service.yaml"
Enter
Sleep 3s

Hide
Type "clear"
Enter
Show

Type "cat examples/split/deployment.yaml"
Enter
Sleep 3s
Show

Hide
Type "clear"
Enter
Show

Type "splinter merge examples/split"
Sleep 750ms
Enter

Sleep 5s