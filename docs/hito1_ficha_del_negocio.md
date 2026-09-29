# Hito 1 · Ficha del negocio

**Pareja:** jandry Paúl Sánchez Murillo · Joseph Damian Loor Chila
**Paralelo:** Aplicación para el Servidor Web A
**Negocio en una línea:** Plataforma de búsqueda de mascotas perdidas con alertas geolocalizadas y gestión de recompensas.

## 1. Negocio de referencia

Plataforma basada en la historia de éxito de *PawBoost* (SaaS de recuperación de mascotas perdidas). La empresa cobra por promocionar avisos de mascotas perdidas mediante campañas automatizadas en redes sociales y alertas locales a una red de voluntarios. Cobra a través de planes de difusión pagados (desde $29 hasta $99 por anuncio impulsado). Según su fundador en Starter Story, declara ingresos superiores a $120,000 mensuales manteniendo un equipo de operaciones mínimo gracias a la automatización de publicaciones y coincidencias.

**Enlace:** https://www.starterstory.com/stories/pawboost

## 2. Caso de contraste

*PetAmberAlert*, una plataforma tradicional de alertas de mascotas por llamadas telefónicas y faxes a clínicas veterinarias. Fracasó y perdió relevancia por no adaptarse a las redes sociales ni al filtrado automático por geolocalización, manteniendo costos operativos muy altos para el volumen de respuesta obtenido.

**Fuente:** https://www.failory.com/blog/pet-finder-failures
**Hipótesis de fondo:** PetAmberAlert dependía de un modelo de costos fijos altos por llamada manual y fax, mientras que PawBoost basa su margen en automatización de APIs publicitarias y comunidades orgánicas, bajando el costo marginal por alerta a casi cero.

## 3. Adaptación al Ecuador

1. **Prevalencia de pagos por transferencia bancaria directa o efectivo:** En Ecuador, el uso de tarjetas de crédito para micro-pagos en línea es bajo. Los usuarios prefieren transferencia o efectivo para el pago de recompensas y publicaciones destacadas.
2. **Alta informalidad en rescatistas y refugios independientes:** Muchas personas que encuentran mascotas son ciudadanos comunes o rescatistas sin personería jurídica que no emiten facturas ni poseen cuentas corporativas.
3. **Poca precisión en direcciones postales:** La búsqueda depende de referencias físicas ("frente al parque central", "junto a la farmacia") más que de códigos postales o coordenadas exacta.

**Qué cambió en el modelo por estas restricciones:** Se agregó el estado `pago_por_verificar` en la entidad `Publicacion` y el campo `ubicacion_referencia` de tipo texto descriptivo. Además, se integró el rol `Comprobador/Agente` para validar manualmente la foto del comprobante de transferencia antes de activar la alerta destacada.

## 4. Modelo de datos

### Entidad: Usuario

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 1 |
| nombre | texto | sí | Jandry Sánchez |
| correo | texto | sí | jandry@ejemplo.com |
| telefono | texto | sí | 0991234567 |
| rol | uno de: dueño, agente | sí | dueño |

### Entidad: ReportePerdido

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 101 |
| usuario_id | referencia a Usuario | sí | 1 |
| nombre_mascota | texto | sí | Firulais |
| especie | uno de: perro, gato, otro | sí | perro |
| raza | texto | sí | Mestizo |
| foto_url | texto | sí | https://img.site/101.jpg |
| ubicacion_referencia | texto | sí | Frente al Parque Central |
| estado | uno de: activo, en_coincidencia, resuelto | sí | activo |
| creado | fecha y hora | sí | 2026-09-28 10:00:00 |

### Entidad: ReporteEncontrado

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 201 |
| usuario_id | referencia a Usuario | sí | 2 |
| especie | uno de: perro, gato, otro | sí | perro |
| foto_url | texto | sí | https://img.site/201.jpg |
| ubicacion | texto | sí | Cerca del mercado municipal |
| estado | uno de: activo, resuelto | sí | activo |
| creado | fecha y hora | sí | 2026-09-28 11:30:00 |

### Entidad: Coincidencia

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 501 |
| reporte_perdido_id | referencia a ReportePerdido | sí | 101 |
| reporte_encontrado_id | referencia a ReporteEncontrado | sí | 201 |
| puntaje_similitud | número decimal | sí | 85.5 |
| estado | uno de: pendiente, confirmado, rechazado | sí | pendiente |
| creado | fecha y hora | sí | 2026-09-28 12:00:00 |

### Relaciones

| Entidades | Cardinalidad | Frase |
|-----------|--------------|-------|
| Usuario — ReportePerdido | 1—N | Un usuario registra muchos reportes de mascotas perdidas |
| Usuario — ReporteEncontrado | 1—N | Un usuario registra muchos reportes de mascotas encontradas |
| ReportePerdido — Coincidencia | 1—N | Un reporte de mascota perdida puede generar varias coincidencias sugeridas |
| ReporteEncontrado — Coincidencia | 1—N | Un reporte de encontrada puede asociarse a varias coincidencias |

### Structs en Go

```go
type Usuario struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nombre   string `gorm:"not null" json:"nombre"`
	Correo   string `gorm:"uniqueIndex;not null" json:"correo"`
	Telefono string `gorm:"not null" json:"telefono"`
	Rol      string `gorm:"not null;default:'dueño'" json:"rol"`
}

type ReportePerdido struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UsuarioID           uint      `gorm:"not null" json:"usuario_id"`
	NombreMascota       string    `gorm:"not null" json:"nombre_mascota"`
	Especie             string    `gorm:"not null" json:"especie"`
	Raza                string    `json:"raza"`
	FotoURL             string    `gorm:"not null" json:"foto_url"`
	UbicacionReferencia string    `gorm:"not null" json:"ubicacion_referencia"`
	Estado              string    `gorm:"not null;default:'activo'" json:"estado"`
	Creado              time.Time `json:"creado"`
}

type ReporteEncontrado struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UsuarioID uint      `gorm:"not null" json:"usuario_id"`
	Especie   string    `gorm:"not null" json:"especie"`
	FotoURL   string    `gorm:"not null" json:"foto_url"`
	Ubicacion string    `gorm:"not null" json:"ubicacion"`
	Estado    string    `gorm:"not null;default:'activo'" json:"estado"`
	Creado    time.Time `json:"creado"`
}

type Coincidencia struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	ReportePerdidoID    uint      `gorm:"not null" json:"reporte_perdido_id"`
	ReporteEncontradoID uint      `gorm:"not null" json:"reporte_encontrado_id"`
	PuntajeSimilitud    float64   `gorm:"type:decimal(5,2)" json:"puntaje_similitud"`
	Estado              string    `gorm:"not null;default:'pendiente'" json:"estado"`
	Creado              time.Time `json:"creado"`
}