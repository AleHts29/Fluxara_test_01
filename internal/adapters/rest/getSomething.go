package rest

// func (h *Handlers) GetProductsAll() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.Background()

// 		vars := mux.Vars(r)
// 		id := vars["id"] // id
// 		fmt.Printf("Esto es id: %+v\n", id)

// 		product, err := h.serviceDb.GetProductsAll(ctx)
// 		if err != nil {
// 			log.Panic("Error en GetProduct")
// 		}

// 		fmt.Printf("Esto es product %+v \n", product)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		json.NewEncoder(w).Encode(product)
// 	}
// }

// func (h *Handlers) GetProduct() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.Background()

// 		vars := mux.Vars(r)
// 		id := vars["id"] // id
// 		fmt.Printf("Esto es id: %+v\n", id)

// 		product, err := h.serviceDb.GetProduct(ctx, id)
// 		if err != nil {
// 			log.Panic("Error en GetProduct")
// 		}

// 		fmt.Printf("Esto es product %+v \n", product)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		json.NewEncoder(w).Encode(product)
// 	}
// }
