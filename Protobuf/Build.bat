del /s /q "..\DouDizhuClient\Assets\Scripts\Network\Proto\*"
del /s /q "..\DouDizhuServer\network\protodef\*"
protoc --proto_path=./ --proto_path=./include --csharp_out=../DouDizhuClient/Assets/Scripts/Network/Proto --go_out=../DouDizhuServer/ *.proto
pause