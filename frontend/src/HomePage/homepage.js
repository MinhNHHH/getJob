import React, { useState } from 'react';
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
} from '@mui/material';

const HomePage = () => {
  const [open, setOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState(null);

  // Sample data - replace with your actual data
  const jobs = [
    {
      id: 1,
      jobTitle: 'Software Engineer',
      companyName: 'Tech Corp',
      location: 'San Francisco, CA',
      description: '<p>Looking for a skilled software engineer to join our team...</p><ul><li>5+ years of experience</li><li>Strong problem-solving skills</li><li>Excellent communication abilities</li></ul>',
    },
    {
      id: 2,
      jobTitle: 'Software Engineer',
      companyName: 'Tech Corp',
      location: 'San Francisco, CA',
      description: '<p>Looking for a skilled softw1111are engineer to join our team...</p><ul><li>5+ years of experience</li><li>Strong problem-solving skills</li><li>Excellent communication abilities</li></ul>',
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

  return (
    <div style={{ padding: '20px' }}>
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
              <TableCell>View Description</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {jobs.map((job) => (
              <TableRow key={job.id}>
                <TableCell>{job.jobTitle}</TableCell>
                <TableCell>{job.companyName}</TableCell>
                <TableCell>{job.location}</TableCell>
                <TableCell>
                  <Button onClick={() => handleClickOpen(job)}>View Description</Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog
        open={open}
        onClose={handleClose}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          {selectedJob?.jobTitle} - {selectedJob?.companyName}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            <div dangerouslySetInnerHTML={{ __html: selectedJob?.description }} />
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Close</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default HomePage;