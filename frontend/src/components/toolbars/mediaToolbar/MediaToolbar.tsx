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
    if (key == "mediaSort" && value == settings.mediaSort) {
      await fetch(`/api/settings`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          id: "mediaSortDirection",
          value:
            settings?.mediaSortDirection.value === "ascending"
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
            system?.scanRunning === "1" ? styles.spinning : styles.svg
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
        system?.metadataRunning == "1" ? (
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
          settingId: "mediaView",
          onClick: () => setSetting("mediaView", "table"),
        },
        {
          text: "Posters",
          settingId: "mediaView",
          id: "posters",
          key: "posters",
          onClick: () => setSetting("mediaView", "posters"),
        },
        {
          text: "Overview",
          settingId: "mediaView",
          id: "overview",
          key: "overview",
          onClick: () => setSetting("mediaView", "overview"),
        },
      ]}
    />,
    <ToolBarItem
      text="Sort"
      index={4}
      key={4}
      settings={settings}
      sortDirection={settings?.mediaSortDirection}
      sort={true}
      icon={<SortIcon className={styles.svg} />}
      selected={selected}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "Monitored/Status",
          id: "monitored/status",
          key: "monitored/status",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "monitored/status"),
        },
        {
          text: "Title",
          id: "title",
          key: "title",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "title"),
        },
        {
          text: "Network",
          id: "network",
          key: "network",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "network"),
        },
        {
          text: "Profile",
          id: "profile",
          key: "profile",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "profile"),
        },
        {
          text: "Episode Count",
          id: "episodes",
          key: "episodes",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "episodes"),
        },
        {
          text: "Size On Disk",
          id: "size",
          key: "size",
          settindId: "mediaSort",
          onClick: () => setSetting("mediaSort", "size"),
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
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "all"),
        },
        {
          text: "Monitored Only",
          id: "monitored",
          key: "monitored",
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "monitored"),
        },
        {
          text: "Unmonitored Only",
          id: "unmonitored",
          key: "unmonitored",
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "unmonitored"),
        },
        {
          text: "Continuing Only",
          id: "continuing",
          key: "continuing",
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "continuing"),
        },
        {
          text: "Ended Only",
          id: "ended",
          key: "ended",
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "ended"),
        },
        {
          text: "Missing Episodes",
          id: "missing",
          key: "missing",
          settindId: "mediaFilter",
          onClick: () => setSetting("mediaFilter", "missing"),
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
