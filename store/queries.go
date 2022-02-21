package store

const (
	insert  = "INSERT INTO customers (name) VALUES (?)"
	get     = "SELECT * FROM customers"
	getByID = "SELECT * FROM customers WHERE id = ?"
	update  = "UPDATE customers SET name = ? WHERE id = ?"
	delete  = "DELETE FROM customers WHERE id = ?"
)
