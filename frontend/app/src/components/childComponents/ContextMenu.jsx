const ContextMenu = ({ items, onItemClick }) => {
  return (
    <div className="context-menu">
      {items.map((item, index) => (
        <div
          key={index}
          className="context-menu-item"
          onClick={() => onItemClick(item)}
        >
          {item.label}
        </div>
      ))}
    </div>
  );
}

export default ContextMenu;