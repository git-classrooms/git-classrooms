$BASE_PATH = $PWD

$SWAGGER_CLIENT_FOLDER = "frontend/src/swagger-client"

docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate -i /local/docs/swagger.yaml -l typescript-axios -o /local/$SWAGGER_CLIENT_FOLDER

Set-Location $SWAGGER_CLIENT_FOLDER

Remove-Item -Recurse -Force .swagger-codegen
@(".gitignore", ".swagger-codegen-ignore", ".npmignore", "git_push.sh", "package.json", "README.md", "tsconfig.json") | ForEach-Object { Remove-Item $PSItem }

Get-ChildItem "apis" -Filter *.ts | Foreach-Object {
    $(
        Write-Output "// @ts-nocheck"
        Get-Content $_.FullName -Raw
    ) | Set-Content $_.FullName
}

Set-Location $BASE_PATH
