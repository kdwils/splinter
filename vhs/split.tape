Output vhs/split.gif

Set Shell "bash"
Set FontSize 20
Set Width 1200
Set Height 600
Set WindowBar Colorful
Set CursorBlink false
Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10

Type "cat examples/merged/merged.yaml"
Enter
Sleep 3s

Hide
Type "clear"
Enter
Show

Type "splinter split examples/merged/merged.yaml -o examples/split -k"
Enter
Sleep 3s
Type "ls examples/split"
Enter
Sleep 3s
Type "cat examples/split/kustomization.yaml"
Enter
Sleep 5s