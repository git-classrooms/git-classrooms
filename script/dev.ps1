
docker compose -f docker-compose.local.yml up -d

$airid=Start-Process -NoNewWindow -PassThru air -c .air.toml

cd frontend
yarn

$yarnid=Start-Process -NoNewWindow -PassThru yarn dev
cd ..

try
{
    while(1) {
        Start-Sleep -Seconds 1
    }
}
finally
{
    $airid | Stop-Process
    $yarnid | Stop-Process
    docker compose -f docker-compose.local.yml stop
}
