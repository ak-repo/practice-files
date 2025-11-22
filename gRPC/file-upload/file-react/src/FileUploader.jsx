import { useState } from "react";

export default function FileUploader() {
  const [file, setFile] = useState(null);

  const upload = async () => {
    const form = new FormData();
    form.append("file", file);

    const res = await fetch("http://localhost:8080/api/v1/files/upload", {
      method: "POST",
      body: form,
    });

    console.log(await res.json());
  };

  return (
    <div className="p-4">
      <input 
        type="file" 
        onChange={(e) => setFile(e.target.files[0])}
      />
      <button 
        onClick={upload} 
        className="bg-blue-600 text-white px-4 py-2 rounded"
      >
        Upload
      </button>
    </div>
  );
}
