# Имя образа
IMAGE_NAME := tg-shortener
# Имя контейнера
CONTAINER_NAME := tg-shortener-container
# Порт (если нужно)
PORT := 8080

.PHONY: build run stop clean logs rebuild

# Собрать Docker-образ
build:
	docker build -t $(IMAGE_NAME) .

# Запустить контейнер в фоновом режиме
run:
	docker run -d --name $(CONTAINER_NAME) $(IMAGE_NAME)
	docker network connect urlshorter_default tg-shortener-container

# Остановить контейнер
stop:
	docker stop $(CONTAINER_NAME)

# Удалить контейнер
rm:
	docker rm $(CONTAINER_NAME)

# Удалить образ
rmi:
	docker rmi $(IMAGE_NAME)

# Просмотр логов
logs:
	docker logs $(CONTAINER_NAME) -f

# Полная очистка (остановить контейнер, удалить его и образ)
clean: stop rm rmi

# Перезапустить контейнер
restart: stop run

rebuild: clean build