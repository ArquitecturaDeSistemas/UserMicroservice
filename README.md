# UserMicroservice
Microservicio de usuario:

- **backend**: 
  - **src**
    - **database**: Conexión a base de datos y donde se establecen las migraciones
    - **middleware**: Funciones relacionadas a la autenticación
    - **models**: Modelos: Inputs que se reciben de la api y modelos que retorna como respuesta.
    - **proto**: Almacena archivos de definición del protocolo gRPC. Contiene las definiciones de los servicios y mensajes que serán utilizados para la comunicación. Solo se codifica el user.proto, lo demás es autogenerado.
    - **repository**: Contiene la lógica de negocio, el cómo se implementan los servicios de usuarios en la bd.
    - **service**: Este servicio implementa la interfaz que se genera del archivo proto
    - **main.go**: Se inicializan los repos, bd y **se define el servidor gRPC** del microservicio.

- **README.md**: Documentación

## Comandos varios

**Comando utilzado para generar los archivos proto (ubicarse en carpeta proto):**
*protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto*

**Correr el programa:**
*go run main.go*
