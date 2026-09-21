# go-semantic-engine-v2
go-semantic-engine-v2


jq -c . get_product_by_SKU > get_product_by_SKU_horizontal
jq -c . get_product_by_SKU

```
                 "get product sku-100"
                          │
                          ▼
                  Intent extraction
                          │
             ┌────────────┴────────────┐
             ▼                         ▼
       structured intent          semantic text
       READ / PRODUCT / SKU       "get product by SKU"
             │                         │
             │                  ┌──────┴──────┐
             │                  ▼             ▼
             │                BM25         Vector
             │                  └──────┬──────┘
             │                         ▼
             └──────────────────► candidate set
                                       │
                                       ▼
                                  reranker
                                       │
                                       ▼
                         inventory.product.get_by_sku
                                       │
                                       ▼
                                sku = "sku-100"

CAPABILITY
inventory.product.get_by_sku
        │
        ├── CAPABILITY_PROTOCOL
        │      endpoint = http://inventory-mcp-server.com:7500/mcp
        │      resource_uri = product://{sku}
        │
        └── CAPABILITY_EXAMPLE
               │
               ├── "Retrieves product information by SKU"
               │       └── EMBEDDING → vector
               │
               ├── "Get product using SKU"
               │       └── EMBEDDING → vector
               │
               ├── "Find product by SKU"
               │       └── EMBEDDING → vector
               │
               └── "Show product details for SKU"
                       └── EMBEDDING → vector


1 | get_by_sku | "Get product by SKU"              | POSITIVE
2 | get_by_sku | "Retrieve product using SKU"     | POSITIVE
3 | get_by_sku | "Find product by SKU"             | POSITIVE

4 | get_by_sku | "Update product by SKU"           | HARD_NEGATIVE
5 | get_by_sku | "Change product stock by SKU"     | HARD_NEGATIVE

query:
"update product sku-100"

examples:

"Get product by SKU"
"Retrieve product using SKU"
"Update product by SKU"
"Change product stock by SKU"                       
```

inventory and product information, product sku information, inventory and product details, product monitoring, inventory monitoring, inventory agent metadata and health, product stock details, inventory available items, product sold, inventory and product stock data, sold and available product figures, update stock level, product price, product name and lead time, change product price, inventory pending itens, stock available pending and sold quantity, create a new product, insert product and inventory, change price and stock data, Updates an inventory inventory quantity sold pending, Updates an existing product price with currency and amount and initial inventory quantity, Retrieves product information for a given SKU, Retrieves general inventory information

Retrieves product information for a given SKU, Inventory and Product details, Product stock data from a given sku, ...

Canonical
Retrieves product information using a product SKU.
Returns product details, inventory information, price,
stock availability, and other product attributes for the specified SKU.

Example Positive
get product by SKU
retrieve product information by SKU
find product using SKU
lookup product details for a SKU
show product information for SKU
get inventory information for a product SKU
retrieve product stock for a SKU
check product details using SKU
find product by product code

Example Negative
update product by SKU
change product information
modify product stock
update product SKU
set product price



@set v1 = '[-0.029558709,0...]'

select :v1

select  ce.id,
		ce."text",
		cer.id,
	    cer."type",
	    c.name,
	    ced.endpoint,
	    ced.uri,
	    e.id,
		e.vector 
from embedding e,
	 capability_example_relation cer,
	 capability c,
	 capability_endpoint ced,
	 capability_example ce
where (e.vector <=> :v1) < 0.8
and e.fk_cap_exp_id = ce.id
and cer.fk_cap_exp = ce.id
and c.id = cer.fk_cap_id 
and c.fk_svc_id = cer.fk_cap_svc_id 
and ced.fk_cap_id = c.id 
and ced.fk_cap_svc_id = c.fk_svc_id
order by (e.vector <=>:v1) asc
limit 10