const sortAndFilter = (
	series: any,
	movies: any,
	profiles: any,
	sort: string,
	sortDirection: string,
	filter: string
) => {
	let filteredMedia: any[] = Object.values(series || {}).concat(
		Object.values(movies || {})
	);
	if (filter == "monitored") {
		filteredMedia = filteredMedia.filter((series: any) => series.monitored);
	} else if (filter == "unmonitored") {
		filteredMedia = filteredMedia.filter((series: any) => !series.monitored);
	} else if (filter == "continuing") {
		filteredMedia = filteredMedia.filter(
			(series: any) => series.status != "Ended"
		);
	} else if (filter == "ended") {
		filteredMedia = filteredMedia.filter(
			(series: any) => series.status == "Ended"
		);
	} else if (filter == "missing") {
		filteredMedia = filteredMedia.filter(
			(series: any) => series.missing_episodes != 0
		);
	}

	let sortedMedia: any[] = filteredMedia;
	const sortSeries = (seriesArray: any[], column: string) => {
		return seriesArray.sort((a: any, b: any) => {
			if (typeof a[column] === "string" && typeof b[column] === "string") {
				return a[column].localeCompare(b[column]);
			} else if (
				typeof a[column] === "number" &&
				typeof b[column] === "number"
			) {
				return a[column] - b[column];
			} else {
				return 0;
			}
		});
	};
	if (sort == "title") {
		sortedMedia = sortSeries(sortedMedia, "id");
	} else if (sort == "monitored/status") {
		sortedMedia = sortSeries(sortedMedia, "monitored");
	} else if (sort == "network") {
		sortedMedia = sortSeries(sortedMedia, "networks");
	} else if (sort == "profile") {
		sortedMedia = sortedMedia.map((series: any) => {
			const profile = profiles[series.profile_id];
			return {
				...series,
				profile_id: profile?.name,
			};
		});
		sortedMedia = sortSeries(sortedMedia, "profile_id");
	} else if (sort == "episodes") {
		sortedMedia = sortSeries(sortedMedia, "episode_count");
	} else if (sort == "size") {
		sortedMedia = sortSeries(sortedMedia, "size");
	}
	if (sortDirection === "descending") {
		sortedMedia = sortedMedia.reverse();
	}
	return sortedMedia;
};

export default sortAndFilter;
