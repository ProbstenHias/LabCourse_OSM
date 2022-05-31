# FaPra OSM 2022
Repository for 'Fachpraktikum Algorithms on OpenStreetMap data' at University Stuttgart by Albi Mema and Matthias Weilinger

# Instructions
### Install Golang
At first, you have to install golang on your System.<br>
Version 1.18 is recommend.<br>
For ubuntu you can follow [this](https://cmatskas.com/install-go-on-wsl-ubuntu-from-the-command-line/) blog post.
### Generate FMI file from PBF file
    1. Navigate into the project folder, so that you are in the directory OSM.
    2. Build the project with command: go build ./src/main.go
    3. Run project with command: ./main {pathToPbfFile}
    4. You can find the resulting fmi file in the same direcory as your pbf file.

### Run Webserver
    1. Navigate into the project folder, so that you are in the directory OSM.
    2. Build the project with command: go build ./src/mainWeb.go
    3. Ron project with command: ./mainWeb {pathToFmiFile} {port}
    4. The GUI can be found at localhost/{port}

### How to use GUI
    - In your browser navigate to localhost/{port}
    - To set a starting point just click anywhere on the map
    - To create a ending point click anywhere on the map
    - Both points will snap to the closest point in water that was created in Task 3
    - Immediately after setting a ending point a shortest path will be drawn on the map
    - When clicking on the marker or on the line more information will be displayed
    - To plan a new route just click anywhere on the map and a new starting point will be set