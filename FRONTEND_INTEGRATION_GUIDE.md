# Guía de Integración para el Frontend - Progress Service

Este documento contiene toda la información necesaria para que el frontend (o una IA de asistencia) implemente la integración con el microservicio de progreso (**progress-service**).

---

## 🚀 Información General
- **Base URL**: `https://gestrym-progress-back.onrender.com/gestrym-progress` (O el host local correspondiente)
- **Auth**: Todos los endpoints requieren un token JWT en el header `Authorization: Bearer <TOKEN>`.

---

## 📊 1. Métricas de Usuario (Peso, Medidas)

### Registrar Métricas
- **Endpoint**: `POST /private/metrics`
- **Body (JSON)**:
```json
{
  "date": "2024-04-24T00:00:00Z",
  "weight": 75.5,
  "height": 180,
  "bodyFat": 15.2,
  "muscleMass": 35.5
}
```

### Obtener Historial de Métricas
- **Endpoint**: `GET /private/metrics/user/:id?limit=10&offset=0`
- **Response**:
```json
{
  "metrics": [
    {
      "id": 1,
      "date": "2024-04-24T00:00:00Z",
      "weight": 75.5,
      "bodyFat": 15.2,
      "muscleMass": 35.5
    }
  ],
  "total": 1
}
```

### Datos para Gráfica de Peso
- **Endpoint**: `GET /private/metrics/user/:id/chart`
- **Response (Cronológico Ascendente)**:
```json
{
  "points": [
    { "date": "2024-03-01T...", "weight": 80.0 },
    { "date": "2024-04-01T...", "weight": 78.5 }
  ]
}
```

---

## 📸 2. Fotos de Progreso

### Subir Foto
- **Endpoint**: `POST /private/photos`
- **Content-Type**: `multipart/form-data`
- **Fields**:
  - `file`: (Archivo de imagen)
  - `type`: "front" | "back" | "side"
  - `date`: "YYYY-MM-DD" (Ej: "2024-04-24")
- **Nota**: El sistema sube la imagen a MinIO automáticamente y retorna éxito.

### Obtener Fotos de un Usuario
- **Endpoint**: `GET /private/photos/user/:id`
- **Response**:
```json
{
  "photos": [
    {
      "id": 1,
      "type": "front",
      "imageUrl": "https://storage.gestrym.com/...",
      "date": "2024-04-24T00:00:00Z"
    }
  ],
  "total": 1
}
```

---

## 🔄 3. Comparación (Antes vs Ahora)

### Obtener Comparativa
- **Endpoint**: `GET /private/comparison/user/:id`
- **Lógica**: Retorna la primera y la última entrada de fotos y métricas.
- **Response**:
```json
{
  "firstPhoto": { "type": "front", "imageUrl": "...", "date": "..." },
  "latestPhoto": { "type": "front", "imageUrl": "...", "date": "..." },
  "firstMetrics": { "weight": 85.0, "bodyFat": 25.0, "date": "..." },
  "latestMetrics": { "weight": 75.0, "bodyFat": 15.0, "date": "..." }
}
```

---

## 🏋️ 4. Progreso por Rutina

### Marcar Rutina Completada
- **Endpoint**: `POST /private/workout-progress`
- **Body (JSON)**:
```json
{
  "workoutId": 123,
  "date": "2024-04-24T15:00:00Z",
  "duration": 60,
  "notes": "Me sentí con mucha energía hoy."
}
```

### Historial de Rutinas
- **Endpoint**: `GET /private/workout-progress/user/:id`
- **Response**:
```json
{
  "progress": [
    {
      "id": 1,
      "workoutId": 123,
      "date": "2024-04-24T15:00:00Z",
      "duration": 60,
      "notes": "..."
    }
  ],
  "total": 1
}
```

---

## 🧑‍🏫 5. Notas del Entrenador

### Crear Nota (Solo Trainers/Admin)
- **Endpoint**: `POST /private/notes`
- **Body (JSON)**:
```json
{
  "userId": 456,
  "message": "Vas por muy buen camino, aumenta 2kg en el press de banca."
}
```

### Ver Notas de un Usuario
- **Endpoint**: `GET /private/notes/user/:id`
- **Response**:
```json
{
  "notes": [
    {
      "id": 1,
      "message": "...",
      "trainerId": 12,
      "date": "2024-04-24T00:00:00Z"
    }
  ],
  "total": 1
}
```

---

## 💡 Instrucciones para la IA del Frontend
1. **Generar Servicios**: Crea un archivo `progressService.ts` usando Axios para realizar estas peticiones.
2. **Tipado**: Define interfaces TypeScript basadas en las respuestas JSON proporcionadas arriba.
3. **Manejo de Formularios**: Para las fotos, asegúrate de usar `FormData`.
4. **Gráficas**: Usa los datos de `/chart` directamente para componentes de Recharts o Chart.js.
5. **Seguridad**: Asegúrate de enviar el token JWT en cada petición.

---

## 🤖 6. Adaptación por Inteligencia Artificial
El backend ahora integra el **ai-service** para adaptar planes de entrenamiento y comida de forma automática.

### 🧠 ¿Cómo funciona?
- **Triggers**: Al llamar a `POST /private/metrics` o `POST /private/photos`, el sistema dispara internamente una señal a la IA.
- **Lógica de Métricas**: Solo se dispara la adaptación si el cambio de peso es significativo (> 0.5kg).
- **Feedback**: Dado que este proceso es **asíncrono** en el backend, no hay cambios inmediatos en la respuesta de la API de progreso. El frontend debe estar preparado para que el usuario reciba notificaciones o vea cambios en sus planes de entrenamiento/comida poco después de subir sus datos.
- **Recomendación**: Puedes mostrar un mensaje al usuario como *"¡Tus nuevos datos están siendo procesados por el Smart Coach para optimizar tu plan!"* al completar el registro.
