import React, { useState } from "react";
import "./App.css";

function App() {
  const [files, setFiles] = useState([]);
  const [tasks, setTasks] = useState({});

  const handleFileChange = (e) => {
    setFiles([...e.target.files]);
  };

  const handleUpload = async () => {
    if (files.length === 0) {
      alert("Please select at least one file.");
      return;
    }

    for (const file of files) {
      const formData = new FormData();
      formData.append("file", file);

      try {
        // Upload the file
        const response = await fetch("http://localhost:8080/api/v1/upload", {
          method: "POST",
          body: formData,
        });

        const data = await response.json();
        if (response.ok) {
          const taskId = data.task_id;
          setTasks((prevTasks) => ({
            ...prevTasks,
            [taskId]: { fileName: file.name, processed: 0, total: 0, status: "" },
          }));
          trackProgress(taskId);
        } else {
          alert(`Error uploading ${file.name}: ${data.error}`);
        }
      } catch (error) {
        alert(`Error uploading ${file.name}`);
        console.error(error);
      }
    }
  };

  const trackProgress = (taskId) => {
    const eventSource = new EventSource(`http://localhost:8080/api/v1/progress-stream/${taskId}`);

    eventSource.onopen = () => {
      console.log("SSE connection opened.");
    };

    eventSource.onmessage = (event) => {
      const progressData = JSON.parse(event.data);
      console.log(progressData)
      setTasks((prevTasks) => ({
        ...prevTasks,
        [taskId]: {
          ...prevTasks[taskId],
          processed: progressData.processed,
          total: progressData.total,
          status: progressData.status,
        },
      }));

      // Close the connection if the task is completed
      if (progressData.status === "completed") {
        eventSource.close();
      }
    };

    eventSource.onerror = (error) => {
      console.error("SSE error:", error);
      eventSource.close();
    };
  };

  return (
    <div className="App">
      <h1>CSV File Uploader</h1>
      <div>
        <input type="file" multiple onChange={handleFileChange} />
        <button onClick={handleUpload}>Upload</button>
      </div>

      <div>
        <h2>Progress</h2>
        {Object.entries(tasks).map(([taskId, task]) => (
          <div key={taskId}>
            <h3>{task.fileName}</h3>
            <p>Task ID: {taskId}</p>
            <p>Processed: {task.processed}</p>
            <p>Total: {task.total}</p>
            <p>Status: {task.status}</p>
            <progress value={task.processed} max={task.total} />
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;