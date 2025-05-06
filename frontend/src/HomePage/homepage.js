import React, { useState, useRef, useEffect } from "react";
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
  CircularProgress,
  TablePagination,
} from "@mui/material";
import api from "../api/fetch";
import axios from "axios";

const HomePage = () => {
  const [open, setOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState(null);
  const [coverLetterOpen, setCoverLetterOpen] = useState(false);
  const [coverLetter, setCoverLetter] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const [isGenerating, setIsGenerating] = useState(false);
  const fileInputRef = useRef(null);

  const [jobs, setJobs] = useState([]);
  const [filters, setFilters] = useState({
    title: "",
    company: "",
    location: "",
  });

  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalJobs, setTotalJobs] = useState(0);

  const fetchJobs = async () => {
    try {
      const response = await api.get(
        `/api/jobs?page=${page + 1}&pageSize=${rowsPerPage}&title=${
          filters.title
        }&company=${filters.company}&location=${filters.location}`
      );
      setJobs(response.data.jobs);
      setTotalJobs(response.data.total_count);
    } catch (error) {
      console.error("Error fetching jobs:", error);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, [page, rowsPerPage, filters.title, filters.company, filters.location]);

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
      setIsGenerating(true);
      // Create FormData instance
      const formData = new FormData();

      // Add job details as JSON string
      formData.append(
        "job_details",
        JSON.stringify({
          job_title: selectedJob.Title,
          company_name: selectedJob.CompanyName,
          location: selectedJob.Location,
          description: selectedJob.Description,
        })
      );

      // Add file if selected
      if (selectedFile) {
        formData.append("resume", selectedFile);
      }

      const response = await api.post("/api/generate-cover-letter", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      if (response) {
        setCoverLetter(response.data);
      } else {
        console.error("No response received");
      }
    } catch (error) {
      console.error("Error generating cover letter:", error);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleFilterChange = (field, value) => {
    setFilters((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  return (
    <div style={{ padding: "20px" }}>
      <Typography variant="h4" gutterBottom>
        Job Listings
      </Typography>

      <Box sx={{ mb: 3 }}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Job Title"
              value={filters.title}
              onChange={(e) => handleFilterChange("title", e.target.value)}
              size="small"
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Company"
              value={filters.company}
              onChange={(e) => handleFilterChange("company", e.target.value)}
              size="small"
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Location"
              value={filters.location}
              onChange={(e) => handleFilterChange("location", e.target.value)}
              size="small"
            />
          </Grid>
        </Grid>
      </Box>

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
            {jobs && jobs.length > 0 ? (
              jobs.map((job) => (
                <TableRow key={job.Id}>
                  <TableCell>{job.Title}</TableCell>
                  <TableCell>{job.CompanyName}</TableCell>
                  <TableCell>{job.Location}</TableCell>
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
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  No jobs found
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={totalJobs}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>

      <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
        <DialogTitle>
          {selectedJob?.Title} - {selectedJob?.CompanyName}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            <div
              dangerouslySetInnerHTML={{ __html: selectedJob?.Description }}
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
          Apply for {selectedJob?.Title} at {selectedJob?.CompanyName}
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
            disabled={isGenerating}
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
                <Button
                  variant="outlined"
                  onClick={generateCoverLetter}
                  disabled={isGenerating}
                  startIcon={
                    isGenerating ? <CircularProgress size={20} /> : null
                  }
                >
                  {isGenerating ? "Generating..." : "Generate Cover Letter"}
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
