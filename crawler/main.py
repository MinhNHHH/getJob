from src.task.base import runner
from src.task.linkedin_crawler import LinkedinCrawler

@runner(task_name="crawl-job")
def main(**kwargs):
    data = []
    crawlers = [LinkedinCrawler(pages=2, job="backend")]
    for crawler in crawlers:
        jobs = crawler.get_jobs()
        data.extend(jobs)
    return data
    
if __name__ == "__main__":
    main()