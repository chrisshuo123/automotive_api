# Automotive APIs

This project mainly showcases my ability in creating Pointer APIs using Golang MVC to show it via HTML Template Frontend and Postman API. For the Postman API, I'll Update how to set this up starting by setting up the MySQL DB.

## Future Release:
- 3D Asset Vehicles per API (Starts from models, but soon can be playable in UNITY 3D. Please take Note: Only Lowpoly Assets) ⏹🚗🚘

## 🚧 Learning Project: [Automotive API]

**Status:** In Progress - Core features completed, optimization ongoing

**What's Working:**
- ✅ Database migration & seeding
- ✅ HTML FrontEnd API
- ✅ RESTful API endpoints

**Currently Learning:**
- 🔄 Update Feature in Postman API
- 🔄 API documentation
- 🔄 Golang Approval Workflow

**Future Improvement:**
- 💡 Add images per vehicle API
- 💡 Add Approval Workflow _(When a user Adds or update new data, the database need to be validated by admin before submitted to MySQL DB)_
- 💡 Add 3D Lowpoly Cars Model _(Why Lowpoly? Because We're a small budgetted team)_

**Tech Stack:** Laravel, MySQL, PHP

# Instructions
In the VSCode IDE on Terminal _(low side of the screen)_, input this:
"git pull https://github.com/chrisshuo123/automotive_api.git main"
If it fails: Download our Repository_

## PART 1: Import Database and Download XAMPP

**Before running the server, please download XAMPP first _(adjust according to your Computer's version: windows, linux, or mac)_ and follow the complete instructions** <br>
<br>
**After finish the installation, open XAMPP, then click "start" on Apache and MySQL, continue to browser and:**
1. input "localhost/phpmyadmin/" on url
2. On the left panel, click "new", give it a name "automotive_api"
3. On automotive_api scheme while selected, click "import" in the top panel
4. Click "choose file", then inside the backend folder, click "automotive_api.sql". _Note: If there's 'toggle Partial Import', other Options such as Foreign Key (FK) Checks are still turn on, please swich the toggle to off._ then, click "import"<br>
<br>
**Return to VSCode, in the terminal run the server _use Golang_, input in the terminal as follow:**
1. "cd backend"
2. "go run main.go"
If correct, it will show "ECHO" with huge text _(DISC: if there's error written in Database, just ignore it)_

**NOTE:** If Golang contains error, then this is the solution. Go to terminal in VSCode:
1. cd backend
2. go mod tidy
_The purpose of "go mod tidy" is to import package '.go' thats essential in running the server through Github Community so Golang can run smoothly_

## PART 2: Postman API Link
If still have no Postman API, please download via the link as follow:
https://www.postman.com/downloads/ 

#### After download, finish installing and open the Postman, proceed by creating a new file,select pointer type _(such as GET or PUT)_, and then copy & paste the link provided below into Postman's URL column:
👉 Pointer GET Cars: 
http://localhost:8000/api/cars
👉 Pointer GET brands:
http://localhost:8000/api/brands
👉 Pointer GET types:
http://localhost:8000/api/types
👉 Pointer PUT Update Cars:
http://localhost:8000/api/cars/_id_

**Disclaimer:** The feature for "Update Cars" still in development 

## PART 3: Website Template: API & CRUD Editor panel
**As long as the server is running, please proceed to _frontend > views_ directory, then:**
1. Right click 'index.html' on the left side bar of the VSCode, click 'Open with live server' to make API Database list appeared in the HTML Website.
2. Right click 'addCar.html' on the left side bar of the VSCode, click 'Open with live server', this is the CRUD Panel where user can create, update, delete API Database via HTML Frontend.






