import styles from "./MediaToolbar.module.scss";
import ToolBarItem from "../../toolBarItem/ToolBarItem";
import Rss from "../../svgs/rss_feed.svg?react";
import Sync from "../../svgs/cached.svg?react";
import AppsIcon from "../../svgs/apps.svg?react";
import ViewIcon from "../../svgs/visibility.svg?react";
import SortIcon from "../../svgs/sort.svg?react";
import FilterIcon from "../../svgs/filter.svg?react";
import TableIcon from "../../svgs/table.svg?react";
import OverviewIcon from "../../svgs/view_list.svg?react";
import LoadingIcon from "../../svgs/loading.svg?react";
import ToolBar from "../../toolBar/ToolBar";
const MediaToolbar = ({
  selected,
  setSelected,
  setContent,
  setIsModalOpen,
  settings,
  system,
  view,
}: any) => {
  const handleOptionsClick = () => {
    setContent(settings);
    setIsModalOpen(true);
  };

  const onRefresh = async () => {
    await fetch(`/api/actions/refresh/metadata`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
  };

  const onUpdate = async () => {
    await fetch(`/api/actions/scan`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
  };

  const setSetting = async (key: string, value: any) => {
    if (key == "media_sort" && value == settings.media_sort.value) {
      await fetch(`/api/settings`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          id: "media_sort_direction",
          value:
            settings?.media_sort_direction.value === "ascending"
              ? "descending"
              : "ascending",
        }),
      });
    }
    await fetch(`/api/settings`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({ id: key, value: value }),
    });
  };

  const leftToolBarItems: any = [
    <ToolBarItem
      text="Scan All"
      index={0}
      key={0}
      icon={
        <Sync
          className={
            system?.scan_running === "1" ? styles.spinning : styles.svg
          }
        />
      }
      onClick={onUpdate}
      selected={selected}
      setSelected={setSelected}
    />,
    <ToolBarItem
      text="Refresh Metadata"
      index={1}
      key={1}
      icon={
        system?.metadata_running == "1" ? (
          <LoadingIcon className={styles.loading} />
        ) : (
          <Rss />
        )
      }
      onClick={onRefresh}
      selected={selected}
      setSelected={setSelected}
    />,
  ];

  const middleToolBarItems: any = [
    <ToolBarItem
      text="Options"
      index={2}
      key={2}
      settings={settings}
      icon={
        view === "table" ? (
          <TableIcon className={styles.svg} />
        ) : view === "posters" ? (
          <AppsIcon className={styles.svg} />
        ) : view === "overview" ? (
          <OverviewIcon className={styles.svg} />
        ) : null
      }
      onClick={handleOptionsClick}
      selected={selected}
      setSelected={setSelected}
    />,
  ];
  const rightToolBarItems: any = [
    <ToolBarItem
      text="View"
      index={3}
      key={3}
      settings={settings}
      icon={<ViewIcon className={styles.svg} />}
      selected={selected}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "Table",
          id: "table",
          key: "table",
          setting_id: "media_view",
          onClick: () => setSetting("media_view", "table"),
        },
        {
          text: "Posters",
          setting_id: "media_view",
          id: "posters",
          key: "posters",
          onClick: () => setSetting("media_view", "posters"),
        },
        {
          text: "Overview",
          setting_id: "media_view",
          id: "overview",
          key: "overview",
          onClick: () => setSetting("media_view", "overview"),
        },
      ]}
    />,
    <ToolBarItem
      text="Sort"
      index={4}
      key={4}
      settings={settings}
      sortDirection={settings?.media_sort_direction.value}
      sort={true}
      icon={<SortIcon className={styles.svg} />}
      selected={selected}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "Monitored/Status",
          id: "monitored/status",
          key: "monitored/status",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "monitored/status"),
        },
        {
          text: "Title",
          id: "title",
          key: "title",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "title"),
        },
        {
          text: "Network",
          id: "network",
          key: "network",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "network"),
        },
        {
          text: "Profile",
          id: "profile",
          key: "profile",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "profile"),
        },
        {
          text: "Episode Count",
          id: "episodes",
          key: "episodes",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "episodes"),
        },
        {
          text: "Size On Disk",
          id: "size",
          key: "size",
          setting_id: "media_sort",
          onClick: () => setSetting("media_sort", "size"),
        },
      ]}
    />,
    <ToolBarItem
      text="Filter"
      index={5}
      key={5}
      icon={<FilterIcon className={styles.svg} />}
      selected={selected}
      settings={settings}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "All",
          id: "all",
          key: "all",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "all"),
        },
        {
          text: "Monitored Only",
          id: "monitored",
          key: "monitored",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "monitored"),
        },
        {
          text: "Unmonitored Only",
          id: "unmonitored",
          key: "unmonitored",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "unmonitored"),
        },
        {
          text: "Continuing Only",
          id: "continuing",
          key: "continuing",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "continuing"),
        },
        {
          text: "Ended Only",
          id: "ended",
          key: "ended",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "ended"),
        },
        {
          text: "Missing Episodes",
          id: "missing",
          key: "missing",
          setting_id: "media_filter",
          onClick: () => setSetting("media_filter", "missing"),
        },
      ]}
    />,
  ];
  return (
    <ToolBar
      leftToolBarItems={leftToolBarItems}
      middleToolBarItems={middleToolBarItems}
      rightToolBarItems={rightToolBarItems}
    />
  );
};

export default MediaToolbar;
