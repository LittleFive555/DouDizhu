Write-Host "=== 开始执行Protocol Buffers编译脚本 ===" -ForegroundColor Cyan
Write-Host ""

# 删除客户端Proto目录中的所有文件
Write-Host "[1/3] 开始清理客户端Proto目录..." -ForegroundColor Yellow
try {
    $clientPath = "../DouDizhuClient/Assets/Scripts/Network/Proto/*"
    if (Test-Path -Path "../DouDizhuClient/Assets/Scripts/Network/Proto") {
        Remove-Item -Path $clientPath -Recurse -Force -ErrorAction Stop
        Write-Host "✓ 客户端Proto目录清理完成" -ForegroundColor Green
    } else {
        Write-Host "⚠ 客户端Proto目录不存在，跳过清理" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ 客户端Proto目录清理失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 删除服务器protodef目录中的所有文件
Write-Host "[2/3] 开始清理服务器protodef目录..." -ForegroundColor Yellow
try {
    $serverPath = "../DouDizhuServer/scripts/network/protodef/*"
    if (Test-Path -Path "../DouDizhuServer/scripts/network/protodef") {
        Remove-Item -Path $serverPath -Recurse -Force -ErrorAction Stop
        Write-Host "✓ 服务器protodef目录清理完成" -ForegroundColor Green
    } else {
        Write-Host "⚠ 服务器protodef目录不存在，跳过清理" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ 服务器protodef目录清理失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 使用protoc编译所有proto文件
Write-Host "[3/3] 开始编译Proto文件..." -ForegroundColor Yellow
try {
    Write-Host "正在调用protoc编译器..." -ForegroundColor White
    # 捕获protoc的输出
    $protocOutput = protoc --proto_path=./ --proto_path=./include --csharp_out="../DouDizhuClient/Assets/Scripts/Network/Proto" --go_out="../DouDizhuServer/scripts/" *.proto 2>&1
    
    # 检查是否有编译错误
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Proto文件编译成功" -ForegroundColor Green
        # 列出编译生成的文件
        Write-Host "生成的客户端文件:" -ForegroundColor White
        Get-ChildItem -Path "../DouDizhuClient/Assets/Scripts/Network/Proto" -Recurse -File | ForEach-Object { 
            Write-Host "  - $($_.Name)" -ForegroundColor Gray
        }
        
        Write-Host "生成的服务器文件:" -ForegroundColor White
        Get-ChildItem -Path "../DouDizhuServer/scripts/" -Recurse -File -Filter "*.pb.go" | ForEach-Object { 
            Write-Host "  - $($_.FullName.Substring((Get-Location).Path.Length + 1))" -ForegroundColor Gray
        }
    } else {
        Write-Host "✗ Proto文件编译失败" -ForegroundColor Red
        Write-Host "错误详情:" -ForegroundColor White
        Write-Host $protocOutput -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ 编译过程中发生异常: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== 脚本执行完成！===" -ForegroundColor Cyan

# 暂停执行，等待用户按任意键继续
Write-Host ""
Write-Host "按任意键继续..." -ForegroundColor White
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")