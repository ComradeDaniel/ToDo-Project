document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('category-form').addEventListener('submit', (e) => {
        e.preventDefault();
        addCategory();
    });

    document.addEventListener('mouseover', (e) => {
        if (e.target.classList.contains('editable-on-hover')) {
            e.target.contentEditable = true;
            e.target.classList.add('editable');
        }
    });

    document.addEventListener('mouseout', (e) => {
        if (e.target.classList.contains('editable-on-hover')) {
            e.target.contentEditable = false;
            e.target.classList.remove('editable');
        }
    });

    let draggedItem = null;

    document.addEventListener('dragstart', (e) => {
        if (e.target.classList.contains('draggable')) {
            draggedItem = e.target;
            setTimeout(() => e.target.classList.add('hidden'), 0);
        }
    });

    document.addEventListener('dragend', (e) => {
        if (draggedItem) {
            draggedItem.classList.remove('hidden');
            draggedItem = null;
        }
    });

    document.addEventListener('dragover', (e) => {
        e.preventDefault();
        if (draggedItem && draggedItem.classList.contains('todo-item')) {
            const closestCategory = getClosestCategory(e.clientY);
            if (closestCategory) {
                const afterElement = getDragAfterElement(closestCategory, e.clientY);
                if (afterElement) {
                    closestCategory.insertBefore(draggedItem, afterElement);
                } else {
                    closestCategory.appendChild(draggedItem);
                }
            }
        } else if (draggedItem && draggedItem.classList.contains('category-container')) {
            const container = document.getElementById('categories');
            const afterElement = getDragAfterElement(container, e.clientY, true);
            if (afterElement) {
                container.insertBefore(draggedItem, afterElement);
            } else {
                container.appendChild(draggedItem);
            }
        }
    });

    function getClosestCategory(y) {
        const categories = [...document.querySelectorAll('.category-container')];
        return categories.reduce((closest, category) => {
            const box = category.getBoundingClientRect();
            const offset = y - box.top - box.height / 2;
            if (offset < 0 && offset > closest.offset) {
                return { offset: offset, element: category };
            } else {
                return closest;
            }
        }, { offset: Number.NEGATIVE_INFINITY }).element;
    }

    function getDragAfterElement(container, y, isCategory = false) {
        const elements = [...container.querySelectorAll(isCategory ? '.category-container:not(.hidden)' : '.draggable:not(.hidden)')];
        return elements.reduce((closest, child) => {
            const box = child.getBoundingClientRect();
            const offset = y - box.top - box.height / 2;
            if (offset < 0 && offset > closest.offset) {
                return { offset: offset, element: child };
            } else {
                return closest;
            }
        }, { offset: Number.NEGATIVE_INFINITY }).element;
    }
});

function addCategory() {
    const categoryTitle = document.getElementById('category-title').value;
    if (categoryTitle.trim() === '') return;

    const categoryContainer = document.createElement('div');
    categoryContainer.className = 'category-container draggable';
    categoryContainer.draggable = true;

    const categoryHeader = document.createElement('div');
    categoryHeader.className = 'category-header';

    const categoryTitleField = document.createElement('div');
    categoryTitleField.className = 'category-title editable-on-hover';
    categoryTitleField.textContent = categoryTitle;

    const deleteCategoryButton = document.createElement('button');
    deleteCategoryButton.className = 'category-delete';
    deleteCategoryButton.textContent = 'Kategorie löschen';
    deleteCategoryButton.onclick = () => categoryContainer.remove();

    categoryHeader.appendChild(categoryTitleField);
    categoryHeader.appendChild(deleteCategoryButton);
    categoryContainer.appendChild(categoryHeader);

    const addButton = document.createElement('button');
    addButton.textContent = 'Aufgabe hinzufügen';
    addButton.onclick = () => addTodoItem(categoryContainer);
    categoryContainer.appendChild(addButton);

    document.getElementById('categories').appendChild(categoryContainer);

    document.getElementById('category-title').value = '';
}

function addTodoItem(categoryContainer) {
    const todoItem = document.createElement('div');
    todoItem.className = 'todo-item draggable';
    todoItem.draggable = true;

    const dragHandle = document.createElement('span');
    dragHandle.className = 'drag-handle';
    dragHandle.innerHTML = '<i class="fas fa-grip-lines"></i>';

    const titleInput = document.createElement('input');
    titleInput.type = 'text';
    titleInput.placeholder = 'Titel';
    titleInput.className = 'editable-on-hover';

    const detailsInput = document.createElement('input');
    detailsInput.type = 'text';
    detailsInput.placeholder = 'Details';
    detailsInput.className = 'editable-on-hover';

    const dueDateInput = document.createElement('input');
    dueDateInput.type = 'datetime-local';
    dueDateInput.className = 'editable-on-hover';

    const completeButton = document.createElement('button');
    completeButton.textContent = 'Erledigt';
    completeButton.className = 'complete';
    completeButton.onclick = () => {
        if (todoItem.style.textDecoration === 'line-through') {
            todoItem.style.textDecoration = 'none';
        } else {
            todoItem.style.textDecoration = 'line-through';
        }
    };

    const deleteButton = document.createElement('button');
    deleteButton.textContent = 'Löschen';
    deleteButton.className = 'delete';
    deleteButton.onclick = () => categoryContainer.removeChild(todoItem);

    todoItem.appendChild(dragHandle);
    todoItem.appendChild(titleInput);
    todoItem.appendChild(detailsInput);
    todoItem.appendChild(dueDateInput);
    todoItem.appendChild(completeButton);
    todoItem.appendChild(deleteButton);

    categoryContainer.appendChild(todoItem);
}
