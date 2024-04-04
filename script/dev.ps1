docker compose -f docker-compose.local.yml up -d

$airid=Start-Process -NoNewWindow -PassThru air -ArgumentList "-c .air.toml" 

cd frontend
yarn

try
{
    yarn dev
}
finally
{
    cd ..
    $airid | Stop-Process
    docker compose -f docker-compose.local.yml stop
}
