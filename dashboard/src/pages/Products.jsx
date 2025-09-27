import React, { useEffect, useState } from "react";
import api from "../api";
import { MdDelete, MdModeEdit } from "react-icons/md";

const Products = () => {
  const [products, setProducts] = useState([]);
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState("All");

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [isEditing, setIsEditing] = useState(false);
  const [currentId, setCurrentId] = useState(null);
  const [showForm, setShowForm] = useState(false);

  const [form, setForm] = useState({
    name: "",
    description: "",
    price: "",
    stock_quantity: "",
    category: "",
    image_url: "",
  });

  // Fetch products
  const fetchProducts = async () => {
    try {
      const res = await api.get("/products"); // public route
      setProducts(res.data.products);
      setFilteredProducts(res.data.products); // initially show all
    } catch (err) {
      console.error(err);
      setError("Failed to load products");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  // Handle input change
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  // Handle add/update product
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const payload = {
        name: form.name,
        description: form.description,
        price: parseFloat(form.price),
        stock_quantity: parseInt(form.stock_quantity, 10),
        category: form.category,
        image_url: form.image_url,
      };

      if (isEditing && currentId) {
        await api.put(`/admin/products/${currentId}`, payload, {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        alert("✅ Product updated successfully");
      } else {
        await api.post("/admin/products", payload, {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        alert("✅ Product created successfully");
      }

      resetForm();
      fetchProducts();
    } catch (err) {
      console.error("Product submit error:", err);
      alert("❌ Failed to save product");
    }
  };

  // Handle edit
  const handleEdit = (product) => {
    setForm({
      name: product.name,
      description: product.description,
      price: product.price,
      stock_quantity: product.stock_quantity,
      category: product.category,
      image_url: product.image_url,
    });
    setIsEditing(true);
    setCurrentId(product.id);
    setShowForm(true);
  };

  // Handle delete
  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this product?"))
      return;
    try {
      await api.delete(`/admin/products/${id}`);
      setProducts(products.filter((p) => p.id !== id));
      setFilteredProducts(filteredProducts.filter((p) => p.id !== id));
    } catch (err) {
      console.error(err);
      alert("Failed to delete product");
    }
  };

  // Reset form
  const resetForm = () => {
    setForm({
      name: "",
      description: "",
      price: "",
      stock_quantity: "",
      category: "",
      image_url: "",
    });
    setIsEditing(false);
    setCurrentId(null);
    setShowForm(false);
  };

  // Handle category filter
  const handleCategoryChange = (e) => {
    const category = e.target.value;
    setSelectedCategory(category);

    if (category === "All") {
      setFilteredProducts(products);
    } else {
      setFilteredProducts(products.filter((p) => p.category === category));
    }
  };

  if (loading) return <p className="p-4">Loading...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold">Products</h1>
        <button
          onClick={() => {
            resetForm();
            setShowForm(true);
          }}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          + Add Product
        </button>
      </div>

      {/* Category Filter */}
      <div className="mb-4">
        <label className="mr-2 font-medium">Filter by Category:</label>
        <select
          value={selectedCategory}
          onChange={handleCategoryChange}
          className="border p-2 rounded"
        >
          <option value="All">All</option>
          {[...new Set(products.map((p) => p.category))].map((cat) => (
            <option key={cat} value={cat}>
              {cat}
            </option>
          ))}
        </select>
      </div>

      {/* Add/Edit Form */}
      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="mb-6 space-y-4 bg-white p-4 shadow rounded"
        >
          <input
            type="text"
            name="name"
            value={form.name}
            onChange={handleChange}
            placeholder="Product Name"
            className="w-full p-2 border rounded"
            required
          />
          <textarea
            name="description"
            value={form.description}
            onChange={handleChange}
            placeholder="Description"
            className="w-full p-2 border rounded"
          />
          <input
            type="number"
            name="price"
            value={form.price}
            onChange={handleChange}
            placeholder="Price"
            className="w-full p-2 border rounded"
            required
          />
          <input
            type="number"
            name="stock_quantity"
            value={form.stock_quantity}
            onChange={handleChange}
            placeholder="Stock Quantity"
            className="w-full p-2 border rounded"
            required
          />
          <input
            type="text"
            name="category"
            value={form.category}
            onChange={handleChange}
            placeholder="Category"
            className="w-full p-2 border rounded"
          />
          <input
            type="text"
            name="image_url"
            value={form.image_url}
            onChange={handleChange}
            placeholder="Image URL"
            className="w-full p-2 border rounded"
          />
          {form.image_url && (
            <img
              src={form.image_url}
              alt="Preview"
              className="w-32 h-32 object-cover rounded border"
            />
          )}
          <div className="space-x-2">
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              {isEditing ? "Update Product" : "Add Product"}
            </button>
            <button
              type="button"
              onClick={resetForm}
              className="px-4 py-2 bg-gray-400 text-white rounded hover:bg-gray-500"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Products Table */}
      <table className="min-w-full border border-gray-200 rounded-lg overflow-hidden bg-white shadow">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-4 py-2 border">ID</th>
            <th className="px-4 py-2 border">Image</th>
            <th className="px-4 py-2 border">Name</th>
            <th className="px-4 py-2 border">Description</th>
            <th className="px-4 py-2 border">Price</th>
            <th className="px-4 py-2 border">Stock</th>
            <th className="px-4 py-2 border">Category</th>
            <th className="px-4 py-2 border">Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredProducts.length > 0 ? (
            filteredProducts.map((p) => (
              <tr key={p.id} className="text-center">
                <td className="px-4 py-2 border">{p.id}</td>
                <td className="px-4 py-2 border">
                  {p.image_url ? (
                    <img
                      src={p.image_url}
                      alt={p.name}
                      className="w-16 h-16 object-cover rounded mx-auto"
                    />
                  ) : (
                    <span className="text-gray-400 italic">No image</span>
                  )}
                </td>
                <td className="px-4 py-2 border">{p.name}</td>
                <td className="px-4 py-2 border">{p.description}</td>
                <td className="px-4 py-2 border">₹{p.price}</td>
                <td className="px-4 py-2 border">{p.stock_quantity}</td>
                <td className="px-4 py-2 border">{p.category}</td>
                <td className="px-4 py-2 border space-x-2">
                  <button
                    onClick={() => handleEdit(p)}
                    className="px-1 py-1 bg-gray-300 rounded-full"
                  >
                    <MdModeEdit />
                  </button>
                  <button
                    onClick={() => handleDelete(p.id)}
                    className="px-1 py-1 bg-gray-300 rounded-full"
                  >
                    <MdDelete />
                  </button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="8" className="text-center py-4 text-gray-500 italic">
                No products found
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Products;
