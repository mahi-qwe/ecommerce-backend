import React, { useEffect, useState } from "react";
import api from "../api";

const Products = () => {
  // State for products list
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // State for showing/hiding forms
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingProductId, setEditingProductId] = useState(null);

  // Form data state
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    price: "",
    stock_quantity: "",
    category: "",
    image_url: "",
  });

  // Clear form data
  const clearForm = () => {
    setFormData({
      name: "",
      description: "",
      price: "",
      stock_quantity: "",
      category: "",
      image_url: "",
    });
  };

  // Fetch all products from API
  const fetchProducts = async () => {
    try {
      setLoading(true);
      const response = await api.get("/products");
      setProducts(response.data.products || []);
      setError("");
    } catch (err) {
      console.error("Error fetching products:", err);
      setError("Failed to load products");
    } finally {
      setLoading(false);
    }
  };

  // Load products when component mounts
  useEffect(() => {
    fetchProducts();
  }, []);

  // Handle input changes in forms
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Add new product
  const handleAddProduct = async (e) => {
    e.preventDefault();
    try {
      await api.post("/admin/products", formData);

      // Reset form and close
      clearForm();
      setShowAddForm(false);

      // Refresh products list
      await fetchProducts();

      setError("");
    } catch (err) {
      console.error("Error adding product:", err);
      setError("Failed to add product");
    }
  };

  // Start editing a product
  const startEditing = (product) => {
    setEditingProductId(product.id);
    setFormData({
      name: product.name || "",
      description: product.description || "",
      price: product.price || "",
      stock_quantity: product.stock_quantity || "",
      category: product.category || "",
      image_url: product.image_url || "",
    });
  };

  // Cancel editing
  const cancelEditing = () => {
    setEditingProductId(null);
    clearForm();
  };

  // Update product
  const handleUpdateProduct = async (e) => {
    e.preventDefault();
    try {
      await api.put(`/admin/products/${editingProductId}`, formData);

      // Reset form and stop editing
      cancelEditing();

      // Refresh products list
      await fetchProducts();

      setError("");
    } catch (err) {
      console.error("Error updating product:", err);
      setError("Failed to update product");
    }
  };

  // Delete product
  const handleDeleteProduct = async (productId, productName) => {
    const confirmDelete = window.confirm(
      `Are you sure you want to delete "${productName}"? This action cannot be undone.`
    );

    if (!confirmDelete) return;

    try {
      await api.delete(`/admin/products/${productId}`);

      // Refresh products list
      await fetchProducts();

      setError("");
    } catch (err) {
      console.error("Error deleting product:", err);
      setError("Failed to delete product");
    }
  };

  // Show loading state
  if (loading) {
    return (
      <div className="p-6">
        <div className="text-center py-8">
          <p className="text-lg text-gray-600">Loading products...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* Page Header */}
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-800">
          Products Management
        </h1>
        <button
          onClick={() => {
            clearForm(); // Clear form before showing
            setShowAddForm(!showAddForm);
          }}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
        >
          {showAddForm ? "Cancel" : "+ Add New Product"}
        </button>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {/* Add Product Form */}
      {showAddForm && (
        <div className="bg-white border border-gray-200 rounded-lg p-6 mb-8 shadow-sm">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">
            Add New Product
          </h2>

          <form onSubmit={handleAddProduct} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Product Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Product Name *
                </label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  placeholder="Enter product name"
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {/* Price */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Price (₹) *
                </label>
                <input
                  type="number"
                  name="price"
                  value={formData.price}
                  onChange={handleInputChange}
                  placeholder="Enter price"
                  required
                  min="0"
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {/* Stock Quantity */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Stock Quantity *
                </label>
                <input
                  type="number"
                  name="stock_quantity"
                  value={formData.stock_quantity}
                  onChange={handleInputChange}
                  placeholder="Enter stock quantity"
                  required
                  min="0"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {/* Category */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Category
                </label>
                <input
                  type="text"
                  name="category"
                  value={formData.category}
                  onChange={handleInputChange}
                  placeholder="Enter category"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Description */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                placeholder="Enter product description"
                rows="3"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-vertical"
              />
            </div>

            {/* Image URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Image URL
              </label>
              <input
                type="url"
                name="image_url"
                value={formData.image_url}
                onChange={handleInputChange}
                placeholder="Enter image URL"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Submit Button */}
            <div className="pt-4">
              <button
                type="submit"
                className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
              >
                Add Product
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Products Table */}
      <div className="bg-white rounded-lg shadow-sm overflow-hidden border border-gray-200">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  ID
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Name
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Description
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Price
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Stock
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Category
                </th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">
                  Image
                </th>
                <th className="px-4 py-3 text-center text-sm font-medium text-gray-700">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {products.length > 0 ? (
                products.map((product) => (
                  <tr key={product.id} className="hover:bg-gray-50">
                    {/* ID */}
                    <td className="px-4 py-4 text-sm text-gray-600">
                      {product.id}
                    </td>

                    {/* Name */}
                    <td className="px-4 py-4">
                      {editingProductId === product.id ? (
                        <input
                          type="text"
                          name="name"
                          value={formData.name}
                          onChange={handleInputChange}
                          className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                          required
                        />
                      ) : (
                        <div className="font-medium text-gray-900">
                          {product.name}
                        </div>
                      )}
                    </td>

                    {/* Description */}
                    <td className="px-4 py-4 max-w-xs">
                      {editingProductId === product.id ? (
                        <textarea
                          name="description"
                          value={formData.description}
                          onChange={handleInputChange}
                          rows="2"
                          className="w-full px-2 py-1 border border-gray-300 rounded text-sm resize-vertical"
                        />
                      ) : (
                        <div className="text-sm text-gray-600">
                          {product.description ? (
                            product.description.length > 60 ? (
                              `${product.description.substring(0, 60)}...`
                            ) : (
                              product.description
                            )
                          ) : (
                            <span className="italic text-gray-400">
                              No description
                            </span>
                          )}
                        </div>
                      )}
                    </td>

                    {/* Price */}
                    <td className="px-4 py-4">
                      {editingProductId === product.id ? (
                        <input
                          type="number"
                          name="price"
                          value={formData.price}
                          onChange={handleInputChange}
                          min="0"
                          step="0.01"
                          className="w-24 px-2 py-1 border border-gray-300 rounded text-sm"
                          required
                        />
                      ) : (
                        <div className="text-sm font-medium text-gray-900">
                          ₹{product.price}
                        </div>
                      )}
                    </td>

                    {/* Stock */}
                    <td className="px-4 py-4">
                      {editingProductId === product.id ? (
                        <input
                          type="number"
                          name="stock_quantity"
                          value={formData.stock_quantity}
                          onChange={handleInputChange}
                          min="0"
                          className="w-20 px-2 py-1 border border-gray-300 rounded text-sm"
                          required
                        />
                      ) : (
                        <div className="text-sm text-gray-600">
                          {product.stock_quantity}
                        </div>
                      )}
                    </td>

                    {/* Category */}
                    <td className="px-4 py-4">
                      {editingProductId === product.id ? (
                        <input
                          type="text"
                          name="category"
                          value={formData.category}
                          onChange={handleInputChange}
                          className="w-32 px-2 py-1 border border-gray-300 rounded text-sm"
                        />
                      ) : (
                        <div className="text-sm text-gray-600">
                          {product.category || (
                            <span className="italic text-gray-400">
                              No category
                            </span>
                          )}
                        </div>
                      )}
                    </td>

                    {/* Image */}
                    <td className="px-4 py-4">
                      {editingProductId === product.id ? (
                        <input
                          type="url"
                          name="image_url"
                          value={formData.image_url}
                          onChange={handleInputChange}
                          placeholder="Image URL"
                          className="w-32 px-2 py-1 border border-gray-300 rounded text-sm"
                        />
                      ) : (
                        <div className="flex flex-col items-center space-y-1">
                          {product.image_url ? (
                            <>
                              <img
                                src={product.image_url}
                                alt={product.name}
                                className="w-12 h-12 object-cover rounded border"
                                onError={(e) => {
                                  e.target.style.display = "none";
                                }}
                              />
                              <a
                                href={product.image_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-xs text-blue-500 hover:text-blue-700 truncate max-w-24"
                                title={product.image_url}
                              >
                                View Image
                              </a>
                            </>
                          ) : (
                            <span className="text-xs text-gray-400 italic">
                              No image
                            </span>
                          )}
                        </div>
                      )}
                    </td>

                    {/* Actions */}
                    <td className="px-4 py-4">
                      <div className="flex justify-center space-x-3">
                        {editingProductId === product.id ? (
                          <>
                            {/* Save Changes Button */}
                            <button
                              onClick={handleUpdateProduct}
                              className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm font-medium transition-colors"
                            >
                              Save
                            </button>
                            {/* Cancel Button */}
                            <button
                              onClick={cancelEditing}
                              className="bg-gray-500 hover:bg-gray-600 text-white px-3 py-1 rounded text-sm font-medium transition-colors"
                            >
                              Cancel
                            </button>
                          </>
                        ) : (
                          <>
                            {/* Edit Button */}
                            <button
                              onClick={() => startEditing(product)}
                              className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-1 rounded text-sm font-medium transition-colors"
                            >
                              Edit
                            </button>
                            {/* Delete Button */}
                            <button
                              onClick={() =>
                                handleDeleteProduct(product.id, product.name)
                              }
                              className="bg-red-500 hover:bg-red-600 text-white px-4 py-1 rounded text-sm font-medium transition-colors"
                            >
                              Delete
                            </button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="8" className="px-4 py-12 text-center">
                    <div className="text-gray-500">
                      <p className="text-lg mb-2">No products found</p>
                      <p className="text-sm">
                        Add your first product to get started!
                      </p>
                    </div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Products;
