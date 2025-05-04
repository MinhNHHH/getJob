import React, { useState, useRef } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  IconButton,
  TextField,
  Box,
  Grid,
} from "@mui/material";
import api from "../api/fetch";

const HomePage = () => {
  const [open, setOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState(null);
  const [coverLetterOpen, setCoverLetterOpen] = useState(false);
  const [coverLetter, setCoverLetter] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const fileInputRef = useRef(null);

  // Sample data - replace with your actual data
  const jobs = [
    {
      id: 1,
      jobTitle: "Software Engineer",
      companyName: "Tech Corp",
      location: "San Francisco, CA",
      description:
        "<p>Looking for a skilled software engineer to join our team...</p><ul><li>5+ years of experience</li><li>Strong problem-solving skills</li><li>Excellent communication abilities</li></ul>",
    },
    {
      id: 2,
      jobTitle: "Software Engineer",
      companyName: "Tech Corp",
      location: "San Francisco, CA",
      description:
        "<p>Looking for a skilled softw1111are engineer to join our team...</p><ul><li>5+ years of experience</li><li>Strong problem-solving skills</li><li>Excellent communication abilities</li></ul>",
    },
    // Add more job listings as needed
  ];

  const handleClickOpen = (job) => {
    setSelectedJob(job);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setSelectedJob(null);
  };

  const handleApply = (job) => {
    setSelectedJob(job);
    setCoverLetterOpen(true);
  };

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current.click();
  };

  const handleCoverLetterSubmit = () => {
    // Here you can implement the logic to submit the cover letter and file
    console.log("Submitting cover letter:", coverLetter);
    if (selectedFile) {
      console.log("Selected file:", selectedFile.name);
      // Add your file upload logic here
    }
    handleCoverLetterClose();
  };

  const handleCoverLetterClose = () => {
    setCoverLetterOpen(false);
    setCoverLetter("");
    setSelectedFile(null);
  };

  const generateCoverLetter = async () => {
    if (!selectedJob) return;

    try {
      // Create FormData instance
      const formData = new FormData();
      
      // Add job details as JSON string
      formData.append('job_details', JSON.stringify({
        job_title: selectedJob.jobTitle,
        company_name: selectedJob.companyName,
        location: selectedJob.location,
        description: selectedJob.description,
      }));

      // Add file if selected
      if (selectedFile) {
        formData.append('resume', selectedFile);
      }

      const response = await api.post("/api/generate-cover-letter", formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });      
      
      if (response) {
        setCoverLetter(response.data);
      } else {
        console.error("No response received");
      }
    } catch (error) {
      console.error("Error generating cover letter:", error);
    }
  };

  return (
    <div style={{ padding: "20px" }}>
      <Typography variant="h4" gutterBottom>
        Job Listings
      </Typography>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="job listings table">
          <TableHead>
            <TableRow>
              <TableCell>Job Title</TableCell>
              <TableCell>Company Name</TableCell>
              <TableCell>Location</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {jobs.map((job) => (
              <TableRow key={job.id}>
                <TableCell>{job.jobTitle}</TableCell>
                <TableCell>{job.companyName}</TableCell>
                <TableCell>{job.location}</TableCell>
                <TableCell>
                  <Button
                    variant="outlined"
                    onClick={() => handleClickOpen(job)}
                    style={{ marginRight: "8px" }}
                  >
                    View Details
                  </Button>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => handleApply(job)}
                  >
                    Apply
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
        <DialogTitle>
          {selectedJob?.jobTitle} - {selectedJob?.companyName}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            <div
              dangerouslySetInnerHTML={{ __html: selectedJob?.description }}
            />
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Cover Letter Dialog */}
      <Dialog
        open={coverLetterOpen}
        onClose={handleCoverLetterClose}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          Apply for {selectedJob?.jobTitle} at {selectedJob?.companyName}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please write your cover letter below:
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="coverLetter"
            label="Cover Letter"
            type="text"
            fullWidth
            multiline
            rows={10}
            variant="outlined"
            value={coverLetter}
            onChange={(e) => setCoverLetter(e.target.value)}
          />
          <Box sx={{ mt: 2 }}>
            <Grid container spacing={2}>
              <Grid item>
                <input
                  type="file"
                  ref={fileInputRef}
                  style={{ display: "none" }}
                  onChange={handleFileChange}
                  accept=".pdf,.doc,.docx"
                />
                <Button variant="outlined" onClick={handleUploadClick}>
                  Upload Resume
                </Button>
              </Grid>
              <Grid item>
                <Button variant="outlined" onClick={generateCoverLetter}>
                  Generate Cover Letter
                </Button>
              </Grid>
            </Grid>
            {selectedFile && (
              <Typography variant="body2" sx={{ mt: 1 }}>
                Selected file: {selectedFile.name}
              </Typography>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCoverLetterClose}>Cancel</Button>
          <Button
            onClick={handleCoverLetterSubmit}
            variant="contained"
            color="primary"
          >
            Submit Application
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default HomePage;
