from src.task.base import JobCrawlBase

class LinkedinCrawler(JobCrawlBase):
    def __init__(self, pages=1, job=''):
        super().__init__(pages, job)

    def get_jobs(self):
        job_ids = self.__get_job_ids()
        jobs = []
        for job_id in job_ids[:2]:
            url_view = f"https://www.linkedin.com/jobs/view/{job_id}"
            print("Crawling:", url_view)
            soup = self.parser_html(url_view)
            if soup:
                title = soup.select_one("h1.top-card-layout__title")
                company = soup.select_one("a.topcard__org-name-link")
                location = soup.select_one("span.topcard__flavor--bullet")
                description = soup.select_one("div.show-more-less-html__markup")
                jobs.append({
                    'title': title.text.strip(),
                    'company_name': company.text.strip(),
                    'company_uri': company.get('href', ''),
                    'location': location.text.strip(),
                    'description': description.prettify(),
                })
        return jobs

    def __get_job_ids(self):
        ids = []
        for page in range(self.pages):
            url = f"https://www.linkedin.com/jobs/search/?keywords={self.job}&start={page*25}"
            soup = self.parser_html(url)
            if not soup:
                return ids
            list_jobs = soup.find_all('div', class_='base-card')
            for job in list_jobs:
                job_id = job.get('data-entity-urn')
                if job_id:
                    parts = job_id.split(':')
                    if len(parts) == 4:
                        ids.append(parts[3])
        return ids

        
if __name__ == "__main__":
    pass