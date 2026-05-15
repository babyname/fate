package repository

import (
	"context"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/poem"
	"github.com/babyname/fate/ent/poemchar"
)

func (m *Repository) QueryPoemsByType(ctx context.Context, poemType poem.Type) ([]*ent.Poem, error) {
	return m.Poem.Query().Where(poem.TypeEQ(poemType)).All(ctx)
}

func (m *Repository) QueryPoemCharsByChar(ctx context.Context, char string) ([]*ent.PoemChar, error) {
	return m.PoemChar.Query().Where(poemchar.CharEQ(char)).WithPoem().All(ctx)
}

func (m *Repository) QueryPoemCharsByPoemID(ctx context.Context, poemID int) ([]*ent.PoemChar, error) {
	return m.PoemChar.Query().Where(poemchar.PoemIDEQ(poemID)).All(ctx)
}

func (m *Repository) InsertPoem(ctx context.Context, title, author, dynasty, content, preface string, keywords, tags []string, poemType poem.Type, source string) (*ent.Poem, error) {
	return m.Poem.Create().
		SetTitle(title).
		SetAuthor(author).
		SetDynasty(dynasty).
		SetContent(content).
		SetPreface(preface).
		SetKeywords(keywords).
		SetTags(tags).
		SetType(poemType).
		SetSource(source).
		Save(ctx)
}

func (m *Repository) InsertPoemChar(ctx context.Context, poemID int, char string, position int, sentence, context string) (*ent.PoemChar, error) {
	return m.PoemChar.Create().
		SetPoemID(poemID).
		SetChar(char).
		SetPosition(position).
		SetSentence(sentence).
		SetContext(context).
		Save(ctx)
}

func (m *Repository) CountPoems(ctx context.Context) (int, error) {
	return m.Poem.Query().Count(ctx)
}

func (m *Repository) CountPoemChars(ctx context.Context) (int, error) {
	return m.PoemChar.Query().Count(ctx)
}
