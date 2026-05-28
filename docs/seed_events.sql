-- Eventos de prueba para desarrollo.
-- Requiere que las tablas ya existan (el backend debe haber corrido AutoMigrate).
-- Uso:
--   docker compose exec -T db mysql -uroot -pmipassword123 eventos_db < docs/seed_events.sql

INSERT INTO events
    (created_at, updated_at, titulo, descripcion, categoria, fecha_hora, duracion, capacidad_total, entradas_vendidas, foto, cancelado)
VALUES
    (NOW(), NOW(), 'Recital de Rock Nacional', 'Una noche de clásicos del rock argentino con bandas en vivo.', 'Musica', '2026-06-15 21:00:00', 180, 1, 0, '', false),
    (NOW(), NOW(), 'Festival de Jazz', 'Tres escenarios con artistas nacionales e internacionales.', 'Musica', '2026-06-20 19:30:00', 240, 300, 0, '', false),
    (NOW(), NOW(), 'Obra de Teatro: Hamlet', 'Clásico de Shakespeare con elenco local.', 'Teatro', '2026-07-02 20:00:00', 150, 120, 0, '', false),
    (NOW(), NOW(), 'Charla TEDx Innovacion', 'Ponencias sobre tecnologia, ciencia y emprendimiento.', 'Conferencia', '2026-07-10 18:00:00', 120, 200, 0, '', false),
    (NOW(), NOW(), 'Partido Amistoso de Futbol', 'Encuentro benefico entre figuras del deporte.', 'Deporte', '2026-06-28 16:00:00', 90, 1000, 0, '', false),
    (NOW(), NOW(), 'Stand Up Comedy Night', 'Lo mejor del humor en vivo. Apto mayores de 16.', 'Humor', '2026-06-12 22:00:00', 100, 80, 0, '', false);
