del /s /q "..\DouDizhuClient\Assets\Scripts\Network\Proto\*"
del /s /q "..\DouDizhuServer\scripts\network\protodef\*"
protoc --proto_path=./ --proto_path=./include --csharp_out=../DouDizhuClient/Assets/Scripts/Network/Proto --go_out=../DouDizhuServer/scripts/ *.proto
pause