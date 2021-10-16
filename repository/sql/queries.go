package sql

const (
	queryWriteBulkNewsTags            = `INSERT INTO news_tags(id, news_id, tag_id, created_at, updated_at) VALUES %s`
	queryReadNewsTags                 = `SELECT news_tags.tag_id, tags.tag FROM news_tags LEFT JOIN tags ON tags.id = news_tags.tag_id WHERE news_tags.news_id = ?`
	queryRemoveNewsTagsByNewsID       = `DELETE FROM news_tags WHERE news_id = ?`
	queryLookupCreateAtNews           = `SELECT id, created_at FROM news WHERE id = ?`
	queryReadNewses                   = `SELECT id, topic_id, title, content, status, created_at, updated_at FROM news ORDER BY created_at DESC`
	queryReadNewsesByStatus           = `SELECT id, topic_id, title, content, status, created_at, updated_at FROM news WHERE status = ? ORDER BY created_at DESC`
	queryReadNewsesByTopicID          = `SELECT id, topic_id, title, content, status, created_at, updated_at FROM news WHERE topic_id = ? AND status = 1 ORDER BY created_at DESC`
	queryReadNewsesByStatusAndTopicID = `SELECT id, topic_id, title, content, status, created_at, updated_at FROM news WHERE topic_id = ? AND status = ? ORDER BY created_at DESC`
	queryWriteNews                    = `INSERT INTO news(id, topic_id, title, content, status, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`
	queryUpdateNews                   = `UPDATE news SET topic_id = ?, title = ?, content = ?, status = ?, created_at = ?, updated_at = ? WHERE id = ?`
	queryRemoveNews                   = `DELETE FROM news WHERE id = ?`
	queryLookupCreateAtTag            = `SELECT id, created_at FROM tags WHERE id = ?`
	queryReadTags                     = `SELECT id, tag, created_at, updated_at FROM tags ORDER BY created_at DESC`
	queryWriteTag                     = `INSERT INTO tags(id, tag, created_at, updated_at) VALUES (?,?,?,?)`
	queryUpdateTag                    = `UPDATE tags SET tag = ?, created_at = ?, updated_at = ? WHERE id = ?`
	queryRemoveTag                    = `DELETE FROM tags WHERE id = ?`
	queryLookupCreateAtTopic          = `SELECT id, created_at FROM topics WHERE id = ?`
	queryReadTopics                   = `SELECT id, title, headline, created_at, updated_at FROM topics ORDER BY created_at DESC`
	queryWriteTopic                   = `INSERT INTO topics(id, title, headline, created_at, updated_at) VALUES (?,?,?,?,?)`
	queryUpdateTopic                  = `UPDATE topics SET title = ?, headline = ?, created_at = ?, updated_at = ? WHERE id = ?`
	queryRemoveTopic                  = `DELETE FROM topics WHERE id = ?`
)
