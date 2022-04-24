# fryazino-temperature-alice
Простой навык для Яндекс Алиса. Озвучивает температуру во Фрязино
docker build --tag frya-temp .
docker run -d --restart=always --name frya-temp -p 13000:3000 -e TZ=Europe/Moscow frya-temp
