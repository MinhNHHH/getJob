from src.task.base import runner
from src.task.linkedin_crawler import LinkedinCrawler
import argparse
import threading


@runner(task_name="crawl-job")
def crawler_worker(crawler, **kwargs):
	try:
		return crawler.get_jobs()
	except Exception as e:
		print(f"Error: {e}")
		return []
        
def crawl_jobs(crawlers):
    """
    Crawl job postings using the given list of crawler classes in parallel using manual thread management.

    Args:
        job (str): The job title to search for.
        pages (int): Number of pages to crawl.
        crawlers (List[Type]): List of crawler classes to use.

    Returns:
        List[Dict]: A combined list of job postings from all crawlers.
    """
    threads = []
    for crawler in crawlers:
        thread = threading.Thread(target=crawler_worker, args=(crawler,))
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()


def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Job Crawler')
    parser.add_argument('--job', type=str, default='', help='Job title to search for')
    parser.add_argument('--pages', type=int, default=1, help='Number of pages to crawl')
    args = parser.parse_args()

    # Print crawling information
    print(f"Crawling for job: {args.job}")
    print(f"Number of pages to crawl: {args.pages}")
    crawlers = [LinkedinCrawler(pages=args.pages, job=args.job)]
    crawl_jobs(crawlers)
    
if __name__ == "__main__":
    main()