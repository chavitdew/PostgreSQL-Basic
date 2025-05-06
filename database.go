package main

func createProduct(product *Product) error {
	_, err := db.Exec(
		"INSERT INTO public.products(name,price) VALUES ($1, $2);",
		product.Name,
		product.Price,
	)
	return err
}

func getProduct(id int) (Product, error) {

	var p Product
	row := db.QueryRow(
		"SELECT id,name,price FROM products WHERE id=$1",
		id,
	)
	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id,name,price FROM products;")

	if err != nil {
		return nil, err
	}

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)

		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	var p Product
	row := db.QueryRow(
		"UPDATE public.products SET name=$2, price=$3 WHERE id=$1 RETURNING id,name,price ;",
		id,
		product.Name,
		product.Price,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM public.products WHERE id=$1;",
		id,
	)
	return err
}
