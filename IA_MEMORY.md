# IA Memory - Implementación de progress-service

Este archivo detalla los cambios y la arquitectura implementada para el microservicio de progreso en el proyecto Gestrym.

## 🏗️ Arquitectura
Se ha seguido estrictamente la **Arquitectura Hexagonal**, separando las responsabilidades en capas:

1.  **Domain**: Interfaces de repositorios y puertos para servicios externos.
2.  **Application**: Casos de uso y DTOs para la transferencia de datos.
3.  **Infrastructure**: Implementaciones GORM de los repositorios y adaptadores para el microservicio de almacenamiento.
4.  **Interfaces/HTTP**: Handlers de Gin para exponer los endpoints.

## 📁 Estructura de Archivos Creados
### Modelos (en `src/common/models`)
- `BodyMetrics.go`: Seguimiento de peso, altura, grasa corporal y masa muscular.
- `ProgressPhoto.go`: Fotos de progreso (frontal, espalda, lateral) con URLs de MinIO.
- `CoachNote.go`: Notas y comentarios de los entrenadores para los usuarios.

### Capa de Dominio (`src/progress/domain`)
- `repositories/`: Interfaces para persistencia.
- `ports/StorageService.go`: Interface para la integración con el servicio de storage.

### Capa de Aplicación (`src/progress/application`)
- `usecases/`: Lógica de negocio (Crear/Obtener métricas, fotos y notas).
- `dtos/`: Estructuras de entrada y salida para la API.

### Capa de Infraestructura (`src/progress/infrastructure`)
- `repositories/`: Implementaciones con GORM.
- `adapters/StorageServiceAdapter.go`: Adaptador que se comunica con el microservicio de storage para subir imágenes.

### Capa de Interfaces (`src/progress/interfaces/http`)
- `handlers/`: Controladores de Gin que manejan las peticiones HTTP.

## 🚀 Endpoints Implementados
Todos los endpoints están bajo el prefijo `/gestrym-progress/private`.

### 📊 Métricas
- `POST /metrics`: Registra nuevas métricas físicas.
- `GET /metrics/user/:id`: Obtiene el historial de métricas de un usuario (soporta paginación).

### 📸 Fotos
- `POST /photos`: Sube una foto de progreso. Acepta `multipart/form-data`.
- `GET /photos/user/:id`: Obtiene las fotos de progreso de un usuario.

### 🧑‍🏫 Notas de Entrenador
- `POST /notes`: Permite a un entrenador dejar una nota a un usuario (Protegido por rol).
- `GET /notes/user/:id`: Obtiene las notas de un usuario.

## 🔐 Reglas de Autorización
- Los usuarios con rol `CLIENTE` solo pueden crear y ver sus propios datos.
- Los entrenadores (`COACH`) pueden crear notas para cualquier usuario y ver su progreso.
- Se utiliza el middleware de JWT existente para extraer el `user_id` y `role_id` del contexto.

## 🛠️ Integración con Storage
Se implementó un adaptador que realiza peticiones POST al microservicio de storage. Las fotos de progreso no se almacenan localmente, solo se guarda la URL retornada por el servicio de almacenamiento (MinIO).

## 📝 Notas Adicionales
- Se han habilitado parámetros de paginación (`limit` y `offset`) en los endpoints GET.
- Se añadió validación para el tipo de foto (`front`, `back`, `side`) mediante tags de Gin.
- Los modelos se registraron en el sistema de migraciones automáticas en `src/common/config/Migrations.go`.

## ✨ Nuevas Funcionalidades (Actualización)
### 🔄 Comparación de Progreso (Antes vs Ahora)
- `GET /comparison/user/:id`: Retorna la primera y la última métrica, junto con la primera y la última foto registrada del usuario. Ideal para visualizaciones de "antes y después".

### 📈 Gráfica de Peso
- `GET /metrics/user/:id/chart`: Retorna un listado cronológico de puntos `(fecha, peso)` listo para ser consumido por bibliotecas de gráficas (como Chart.js o Recharts).

### 🏋️ Registro de Progreso por Rutina
- `POST /workout-progress`: Permite marcar una rutina como completada, registrando la fecha, duración y notas adicionales.
- `GET /workout-progress/user/:id`: Obtiene el historial de rutinas completadas por el usuario.
- **Integración**: Este endpoint utiliza el ID de la rutina proveniente del `training-service` para vincular el progreso.

## 🤖 Integración con AI Service (Adaptación Inteligente)
Se ha implementado una integración reactiva con el microservicio de IA para permitir que los planes de entrenamiento y nutrición se adapten automáticamente al progreso del usuario.

### ⚙️ Implementación Técnica
- **Port/Interface**: `AIService` en la capa de dominio define los métodos `AdaptTraining` y `AdaptNutrition`.
- **Adaptador**: `AIServiceAdapter` implementa la comunicación HTTP con el microservicio de IA.
- **Inyección**: El servicio se inyecta en los casos de uso de métricas corporales y fotos de progreso.

### 🧠 Lógica de Activación (Triggers)
1.  **Creación de Métricas**: Al registrar peso o medidas, el sistema evalúa si el cambio es significativo (umbral de **> 0.5kg**). Si se cumple, se solicita una re-evaluación de los planes.
2.  **Fotos de Progreso**: Subir una nueva foto siempre dispara una solicitud de adaptación, ya que puede representar cambios visuales no capturados por las básculas.

### 🚀 Optimización y Resiliencia
- **Ejecución Asíncrona**: Las llamadas a la IA se ejecutan en goroutines separadas. El usuario recibe su confirmación de guardado inmediatamente sin esperar la respuesta de la IA.
- **Tolerancia a Fallos**: Si el microservicio de IA no está disponible o falla, el error se registra en logs pero el flujo principal del `progress-service` continúa exitosamente.
- **Debounce Natural**: Al validar contra el registro anterior, evitamos llamadas redundantes si los cambios son mínimos.
