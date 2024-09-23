import styles from "./MassEditorToolbar.module.scss";
import ToolBar from "../../toolBar/ToolBar";
import ToolBarItem from "../../toolBarItem/ToolBarItem";
import SortIcon from "../../svgs/sort.svg?react";
import FilterIcon from "../../svgs/filter.svg?react";
const MassEditorToolbar = ({ selected, setSelected, settings }: any) => {
  const setSetting = async (key: string, value: any) => {
    if (key == "massEditorSort" && value == settings.massEditorSort) {
      await fetch(`/api/settings/${key}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          id: "massEditorSortDirection",
          value:
            settings?.massEditorSortDirection == "ascending"
              ? "descending"
              : "ascending",
        }),
      });
    }
    await fetch(`/api/settings/${key}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({ id: key, value: value }),
    });
  };

  const rightToolBarItems: any = [
    <ToolBarItem
      text="Sort"
      index={4}
      key={4}
      settings={settings}
      sortDirection={settings?.massEditorSortDirection}
      sort={true}
      icon={<SortIcon className={styles.svg} />}
      selected={selected}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "Monitored/Status",
          id: "monitored/status",
          key: "monitored/status",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "monitored/status"),
        },
        {
          text: "Title",
          id: "title",
          key: "title",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "title"),
        },
        {
          text: "Network",
          id: "network",
          key: "network",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "network"),
        },
        {
          text: "Profile",
          id: "profile",
          key: "profile",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "profile"),
        },
        {
          text: "Episode Count",
          id: "episodes",
          key: "episodes",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "episodes"),
        },
        {
          text: "Size On Disk",
          id: "size",
          key: "size",
          settingId: "massEditorSort",
          onClick: () => setSetting("massEditorSort", "size"),
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
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "all"),
        },
        {
          text: "Monitored Only",
          id: "monitored",
          key: "monitored",
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "monitored"),
        },
        {
          text: "Unmonitored Only",
          id: "unmonitored",
          key: "unmonitored",
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "unmonitored"),
        },
        {
          text: "Continuing Only",
          id: "continuing",
          key: "continuing",
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "continuing"),
        },
        {
          text: "Ended Only",
          id: "ended",
          key: "ended",
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "ended"),
        },
        {
          text: "Missing Episodes",
          id: "missing",
          key: "missing",
          settingId: "massEditorFilter",
          onClick: () => setSetting("massEditorFilter", "missing"),
        },
      ]}
    />,
  ];
  return <ToolBar rightToolBarItems={rightToolBarItems} />;
};
export default MassEditorToolbar;
