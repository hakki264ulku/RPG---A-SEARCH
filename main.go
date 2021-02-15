package main

import (
	"bufio"
	"os"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/BiesGo/sdlWorkSpace/rpg/ui"
)

func main() {

	//game.GenerateRandomMaze(30, 30)

	file, err := os.Open("./ui/assets/mainMenuInfo.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
	}
	ui := &ui.UI2d{}
	game.Run(ui)

}

//game currently starts in Main Menu
// Start button leads user to Selection Menu
// in selection menu
//		start button leads user to end menu
//		edit button leads user to edit menu

//IN LEVEL SELECTION MENU
//after selecting a level if show preview box checked a preview appears
//user can use arrow keys(left,right,up,down) to traverse maze
//user can hold lshift for minimap
//while holding lshift pageUp & pageDown buttons zooms in & out respectively

// IN EDIT MODE
//for tool menu hold "tab" or "lshift"
//click on tool menu choices or use numbers 1,2,3...
//start button in tool menu leads user to end menu
// ps. textures are created randomly within their scope
//left click place(if possible such as wall and floor uses same layer but door is 1 layer above)
//right click remove most upper layer
// "s" hard save
// "Backspace" reload latest save from file

//todo list
//code beauty(there may be unused/unnecessary parts left while changing whole structure)
//auto maze builder within scope
//left right up down buttons for travelling on level(currently left upmost part is shown automaticly)
// entity collision detection in build menu ie. doors cant overlap & there cant be more than 1 main character
// ~H~
