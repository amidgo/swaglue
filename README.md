## SWAGLUE Инструмент для раздельного написания Swagger 3.0

# Основная концепция
Раздельное написание компонентов OpenApi, а имено: схем ("#/components/schemas/DTO"), запросов, схем авторизации

# Usage
```bash
-components string
      components with name and directory path, example --components=<name>=<path>,<name>=<path>
-head string
      head swagger yaml, should be in .yaml format
-output string
      output file
-paths string
      path to paths directory
```
# Пример использования
Базовый пример использования ref для запроса

Наш основной файл, swagger.yaml
```yaml 
# Путь запроса
/teachers:
# Метод
  get:
# Референс, необходим для нахождения нужного файла
    $ref: "#/paths/teachers/get"
```
Для того чтобы программа подставила нужный нам Yaml в данный метод, необходимо в создать папке, которую мы указали во флаге -paths см. Usage, папку teachers и файл get.yaml в ней

Пример содержания файла get.yaml

```yaml 
tags:
  - teachers
summary: Get All users who are teachers with pagination
security:
  - cookieAuth: []
parameters:
  - name: offset
    in: query
    required: true
    schema:
      $ref: "#/components/schemas/Offset"
  - name: limit
    in: query
    required: true
    schema:
      $ref: "#/components/schemas/Limit"
responses:
  200:
    description: ""
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/UserPageData"
```

при запуске утилиты она спарсит все $ref: "#/paths/<Путь к файлу с запросом>"
И заменит их содержанием этих самых файлов

Итоговый результат:
```yaml
    /teachers:
        get:
            tags:
                - teachers
            summary: Get All users who are teachers with pagination
            security:
                - cookieAuth: []
            parameters:
                - name: offset
                  in: query
                  required: true
                  schema:
                    $ref: "#/components/schemas/Offset"
                - name: limit
                  in: query
                  required: true
                  schema:
                    $ref: "#/components/schemas/Limit"
            responses:
                200:
                    description: ""
                    content:
                        application/json:
                            schema:
                                $ref: "#/components/schemas/UserPageData"
```

Все компоненты парсятся автоматически 
Пример использования

Начальный файл (head) swagger.yaml

```yaml
openapi: 3.0.0
info:
  title: ESUS API
  description: esus
  termsOfService: http://swagger.io/terms/
  license:
    name: Apache 2.0
    url: http://www.apache.org/licenses/LICENSE-2.0.html
  version: 1.0.0
externalDocs:
  description: POFIG NA DESCRIPTION
  url: http://swagger.io
servers:
  - url: http://localhost:10024/api/v1
tags:
...
paths:
...
components:

```

Для того чтобы добавить состовляющую в components необходимо добавить флаг --components=component_name=component_content_path

Например
--components schemas=./schemas
Данный флаг сигнализирует программе о том что мы 
должны добавить значение schemas к components,
должны достать все .yaml файлы из папки schemas и вставить их содержимое в schemas

содержание файла ./schemas/users/User.yaml

```yaml
type: object
properties:
  id:
    $ref: "#/components/schemas/ID32"
  name:
    $ref: "#/components/schemas/UserName"
  fathername:
    $ref: "#/components/schemas/UserFatherName"
  surname:
    $ref: "#/components/schemas/UserSurname"
```
Итоговый результат

```yaml
openapi: 3.0.0
info:
  title: ESUS API
  description: esus
  termsOfService: http://swagger.io/terms/
  license:
    name: Apache 2.0
    url: http://www.apache.org/licenses/LICENSE-2.0.html
  version: 1.0.0
externalDocs:
  description: POFIG NA DESCRIPTION
  url: http://swagger.io
servers:
  - url: http://localhost:10024/api/v1
tags:
...
paths:
...
components:
    schemas:
        User:
            type: object
            properties:
              id:
                $ref: "#/components/schemas/ID32"
              name:
                $ref: "#/components/schemas/UserName"
              fathername:
                $ref: "#/components/schemas/UserFatherName"
              surname:
                $ref: "#/components/schemas/UserSurname"
```

Таким образом мы можем грубо говоря склеивать наш сваггер из маленьких кусочков, которые удобны к изменению