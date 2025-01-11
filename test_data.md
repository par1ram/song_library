#### Таблица `groups`

```sql
INSERT INTO groups (group_name) VALUES
('The Beatles'),
('Queen'),
('Led Zeppelin');
```

#### Таблица `songs`

```sql
INSERT INTO songs (song_name, release_date, text, link, group_id) VALUES
('Hey Jude', '1968-08-26', 'Hey Jude, don''t make it bad...', 'https://example.com/heyjude', (SELECT id FROM groups WHERE group_name = 'The Beatles')),
('Bohemian Rhapsody', '1975-10-31', 'Is this the real life? Is this just fantasy...', 'https://example.com/bohemianrhapsody', (SELECT id FROM groups WHERE group_name = 'Queen')),
('Stairway to Heaven', '1971-11-08', 'There''s a lady who''s sure all that glitters is gold...', 'https://example.com/stairwaytoheaven', (SELECT id FROM groups WHERE group_name = 'Led Zeppelin'));
```

**Запрос**: `POST /songs/filter`

```json
{
	"group": "The Beatles",
	"song": "Hey Jude",
	"release_date": "1968-08-26",
	"limit": 10,
	"offset": 0
}
```

**Запрос**: PUT /songs/update

```json
{
	"id": 1,
	"group_id": 1,
	"song_name": "Hey Jude (Remastered)",
	"text": "Hey Jude, don't make it bad...",
	"release_date": "1968-08-26",
	"link": "https://example.com/heyjude-remastered"
}
```

**Запрос**: `PATCH /songs/patch`

```json
{
	"id": 1,
	"song_name": "Hey Jude (Live)"
}
```
