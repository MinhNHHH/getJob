from src.task.base import JobCrawlBase
from concurrent.futures import ThreadPoolExecutor, as_completed

class LinkedinCrawler(JobCrawlBase):
    def __init__(self, pages=1, job=''):
        super().__init__(pages, job)

    def get_jobs(self):
        job_ids = self.__get_job_ids()
        jobs = []
        max_workers = 2
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            future_to_job_id = {executor.submit(self.crawl_job, job_id): job_id for job_id in job_ids}
            for future in as_completed(future_to_job_id):
                result = future.result()
                if result:
                    jobs.append(result)
        return jobs

    def crawl_job(self,job_id):
        url_view = f"https://www.linkedin.com/jobs/view/{job_id}"
        print("Crawling:", url_view)
        soup = self.parser_html(url_view)
        if soup:
            title = soup.select_one("h1.top-card-layout__title")
            company = soup.select_one("a.topcard__org-name-link")
            location = soup.select_one("span.topcard__flavor--bullet")
            description = soup.select_one("div.show-more-less-html__markup")
            return {
                'title': title.text.strip(),
                'company_name': company.text.strip(),
                'company_uri': company.get('href', ''),
                'location': location.text.strip(),
                'description': description.prettify(),
            }
            
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