package models

type Book struct {
	ID       int     `json:"id"`
	Title    string  `json:"title" db:"title"`
	NumPages int     `json:"numPages" db:"num_pages"`
	Author   string  `json:"author" db:"author"`
	Rating   float64 `json:"rating" db:"rating"`
}

func (b *Book) Validate() map[string]string {
	errors := make(map[string]string)

	if b.Title == "" {
		errors["title"] = "Title is required"
	}

	if b.Author == "" {
		errors["author"] = "Author is required"
	}

	if b.NumPages < 0 {
		errors["numPages"] = "Number of pages cannot be negative"
	}

	if b.Rating < 0 || b.Rating > 5 {
		errors["rating"] = "Rating must be between 0 and 5"
	}

	return errors
}
