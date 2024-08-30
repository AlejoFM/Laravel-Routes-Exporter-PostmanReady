
# Laravel Routes for Postman

With the .exe file it will generate an entire collection of all the routes available from you api.php file in your Laravel project.




## Installation

To build the .exe file, first you would need to have Go available in your device.

```bash
  https://go.dev/dl/
  Choose your OS, execute the installer and after finished installing. 
  Open your terminal, and write : go version
  If something like " go version go1.23.0 windows/amd64 " appears, it means go installed
  succesfully and now its ready to use! 🐭
```
After installing GO :
```
  Go to the repository path and run the following command in your terminal : 
  go build -o {executable-name}.exe ./cmd/generate
```
  This will generate a .exe file, move that file to the root of any Laravel project and run it. It will generate a postman-collection.json file ready to import in your postman

